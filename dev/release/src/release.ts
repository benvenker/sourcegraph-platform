import { readFileSync, rmdirSync, writeFileSync } from 'fs'
import * as path from 'path'

import commandExists from 'command-exists'
import { addMinutes } from 'date-fns'

import * as batchChanges from './batchChanges'
import * as changelog from './changelog'
import { Config, releaseVersions } from './config'
import {
    getAuthenticatedGitHubClient,
    listIssues,
    getTrackingIssue,
    createChangesets,
    CreatedChangeset,
    createTag,
    ensureTrackingIssues,
    releaseName,
    commentOnIssue,
    queryIssues,
    IssueLabel,
    createLatestRelease,
} from './github'
import { ensureEvent, getClient, EventOptions, calendarTime } from './google-calendar'
import { postMessage, slackURL } from './slack'
import {
    cacheFolder,
    formatDate,
    timezoneLink,
    hubSpotFeedbackFormStub,
    ensureDocker,
    changelogURL,
    ensureReleaseBranchUpToDate,
} from './util'

const sed = process.platform === 'linux' ? 'sed' : 'gsed'

export type StepID =
    | 'help'
    // release tracking
    | 'tracking:timeline'
    | 'tracking:issues'
    // branch cut
    | 'changelog:cut'
    // release
    | 'release:status'
    | 'release:create-candidate'
    | 'release:stage'
    | 'release:add-to-batch-change'
    | 'release:finalize'
    | 'release:close'
    // util
    | 'util:clear-cache'
    // testing
    | '_test:google-calendar'
    | '_test:slack'
    | '_test:batchchange-create-from-changes'
    | '_test:config'
    | '_test:dockerensure'

/**
 * Runs given release step with the provided configuration and arguments.
 */
export async function runStep(config: Config, step: StepID, ...args: string[]): Promise<void> {
    if (!steps.map(({ id }) => id as string).includes(step)) {
        throw new Error(`Unrecognized step ${JSON.stringify(step)}`)
    }
    await Promise.all(
        steps
            .filter(({ id }) => id === step)
            .map(async step => {
                if (step.run) {
                    await step.run(config, ...args)
                }
            })
    )
}

interface Step {
    id: StepID
    description: string
    run?: ((config: Config, ...args: string[]) => Promise<void>) | ((config: Config, ...args: string[]) => void)
    argNames?: string[]
}

const steps: Step[] = [
    {
        id: 'help',
        description: 'Output help text about this tool',
        argNames: ['all'],
        run: (_config, all) => {
            console.error('Sourcegraph release tool - https://handbook.sourcegraph.com/engineering/releases')
            console.error('\nUSAGE\n')
            console.error('\tyarn run release <step>')
            console.error('\nAVAILABLE STEPS\n')
            console.error(
                steps
                    .filter(({ id }) => all || !id.startsWith('_'))
                    .map(
                        ({ id, argNames, description }) =>
                            '\t' +
                            id +
                            (argNames && argNames.length > 0
                                ? ' ' + argNames.map(argumentName => `<${argumentName}>`).join(' ')
                                : '') +
                            '\n\t\t' +
                            description
                    )
                    .join('\n') + '\n'
            )
        },
    },
    {
        id: 'tracking:timeline',
        description: 'Generate a set of Google Calendar events for a MAJOR.MINOR release',
        run: async config => {
            const { upcoming: release } = await releaseVersions(config)
            const name = releaseName(release)
            const events: EventOptions[] = [
                {
                    title: `Security Team to Review Release Container Image Scans ${name}`,
                    description: '(This is not an actual event to attend, just a calendar marker.)',
                    anyoneCanAddSelf: true,
                    attendees: [config.teamEmail],
                    transparency: 'transparent',
                    ...calendarTime(config.oneWorkingWeekBeforeRelease),
                },
                {
                    title: `Cut Sourcegraph ${name}`,
                    description: '(This is not an actual event to attend, just a calendar marker.)',
                    anyoneCanAddSelf: true,
                    attendees: [config.teamEmail],
                    transparency: 'transparent',
                    ...calendarTime(config.oneWorkingDayBeforeRelease),
                },
                {
                    title: `Release Sourcegraph ${name}`,
                    description: '(This is not an actual event to attend, just a calendar marker.)',
                    anyoneCanAddSelf: true,
                    attendees: [config.teamEmail],
                    transparency: 'transparent',
                    ...calendarTime(config.releaseDate),
                },
                {
                    title: `Deploy Sourcegraph ${name} to managed instances`,
                    description: '(This is not an actual event to attend, just a calendar marker.)',
                    anyoneCanAddSelf: true,
                    attendees: [config.teamEmail],
                    transparency: 'transparent',
                    ...calendarTime(config.oneWorkingDayAfterRelease),
                },
            ]

            if (!config.dryRun.calendar) {
                const googleCalendar = await getClient()
                for (const event of events) {
                    console.log(`Create calendar event: ${event.title}: ${event.startDateTime || 'undefined'}`)
                    await ensureEvent(event, googleCalendar)
                }
            } else {
                console.log('dryRun.calendar=true, skipping calendar event creation', events)
            }
        },
    },
    {
        id: 'tracking:issues',
        description: 'Generate GitHub tracking issue for the configured release',
        run: async (config: Config) => {
            const {
                releaseDate,
                captainGitHubUsername,
                oneWorkingDayBeforeRelease,
                oneWorkingDayAfterRelease,
                captainSlackUsername,
                slackAnnounceChannel,
                dryRun,
            } = config
            const { upcoming: release } = await releaseVersions(config)
            const date = new Date(releaseDate)

            // Create issue
            const trackingIssues = await ensureTrackingIssues({
                version: release,
                assignees: [captainGitHubUsername],
                releaseDate: date,
                oneWorkingDayBeforeRelease: new Date(oneWorkingDayBeforeRelease),
                oneWorkingDayAfterRelease: new Date(oneWorkingDayAfterRelease),
                dryRun: dryRun.trackingIssues || false,
            })
            console.log('Rendered tracking issues', trackingIssues)

            // If at least one issue was created, post to Slack
            if (trackingIssues.find(({ created }) => created)) {
                const name = releaseName(release)
                const releaseDateString = slackURL(formatDate(date), timezoneLink(date, `${name} release`))
                let annoncement = `:mega: *${name} release*

:captain: Release captain: @${captainSlackUsername}
:spiral_calendar_pad: Scheduled for: ${releaseDateString}
:pencil: Tracking issues:
${trackingIssues.map(index => `- ${slackURL(index.title, index.url)}`).join('\n')}`
                if (release.patch !== 0) {
                    const patchRequestTemplate = `https://github.com/sourcegraph/sourcegraph/issues/new?assignees=&labels=team%2Fdistribution&template=request_patch_release.md&title=${release.version}%3A+`
                    annoncement += `\n\nIf you have changes that should go into this patch release, ${slackURL(
                        'please *file a patch request issue*',
                        patchRequestTemplate
                    )}, or it will not be included.`
                }
                if (!dryRun.slack) {
                    await postMessage(annoncement, slackAnnounceChannel)
                    console.log(`Posted to Slack channel ${slackAnnounceChannel}`)
                }
            } else {
                console.log('No tracking issues were created, skipping Slack announcement')
            }
        },
    },
    {
        id: 'changelog:cut',
        description: 'Generate pull requests to perform a changelog cut for branch cut',
        argNames: ['changelogFile'],
        run: async (config, changelogFile = 'CHANGELOG.md') => {
            const { upcoming: release } = await releaseVersions(config)
            const prMessage = `changelog: cut sourcegraph@${release.version}`
            await createChangesets({
                requiredCommands: [],
                changes: [
                    {
                        owner: 'sourcegraph',
                        repo: 'sourcegraph',
                        base: 'main',
                        head: `changelog-${release.version}`,
                        title: prMessage,
                        commitMessage: prMessage,
                        edits: [
                            (directory: string) => {
                                console.log(`Updating '${changelogFile} for ${release.format()}'`)
                                const changelogPath = path.join(directory, changelogFile)
                                let changelogContents = readFileSync(changelogPath).toString()

                                // Convert 'unreleased' to a release
                                const releaseHeader = `## ${release.format()}`
                                const unreleasedHeader = '## Unreleased'
                                changelogContents = changelogContents.replace(unreleasedHeader, releaseHeader)

                                // Add a blank changelog template for the next release
                                changelogContents = changelogContents.replace(
                                    changelog.divider,
                                    changelog.releaseTemplate
                                )

                                // Update changelog
                                writeFileSync(changelogPath, changelogContents)
                            },
                        ],
                    },
                ],
                dryRun: config.dryRun.changesets,
            })
        },
    },
    {
        id: 'release:status',
        description: 'Post a message in Slack summarizing the progress of a release',
        run: async config => {
            const githubClient = await getAuthenticatedGitHubClient()
            const { upcoming: release } = await releaseVersions(config)

            const trackingIssue = await getTrackingIssue(githubClient, release)
            if (!trackingIssue) {
                throw new Error(`Tracking issue for version ${release.version} not found - has it been created yet?`)
            }

            const blockingQuery = 'is:open org:sourcegraph label:release-blocker'
            const blockingIssues = await listIssues(githubClient, blockingQuery)
            const blockingIssuesURL = `https://github.com/issues?q=${encodeURIComponent(blockingQuery)}`
            const blockingMessage =
                blockingIssues.length === 0
                    ? 'There are no release-blocking issues'
                    : `There ${
                          blockingIssues.length === 1
                              ? 'is 1 release-blocking issue'
                              : `are ${blockingIssues.length} release-blocking issues`
                      }`

            const message = `:mega: *${release.version} Release Status Update*

* Tracking issue: ${trackingIssue.url}
* ${blockingMessage}: ${blockingIssuesURL}`
            if (!config.dryRun.slack) {
                await postMessage(message, config.slackAnnounceChannel)
            }
        },
    },
    {
        id: 'release:create-candidate',
        description: 'Generate the Nth release candidate. Set <candidate> to "final" to generate a final release',
        argNames: ['candidate'],
        run: async (config, candidate) => {
            if (!candidate) {
                throw new Error('Candidate information is required (either "final" or a number)')
            }
            const { upcoming: release } = await releaseVersions(config)
            const branch = `${release.major}.${release.minor}`
            const tag = `v${release.version}${candidate === 'final' ? '' : `-rc.${candidate}`}`
            ensureReleaseBranchUpToDate(branch)
            await createTag(
                await getAuthenticatedGitHubClient(),
                {
                    owner: 'sourcegraph',
                    repo: 'sourcegraph',
                    branch,
                    tag,
                },
                config.dryRun.tags || false
            )
        },
    },
    {
        id: 'release:stage',
        description: 'Open pull requests and a batch change staging a release',
        run: async config => {
            const { slackAnnounceChannel, dryRun } = config
            const { upcoming: release, previous } = await releaseVersions(config)
            // ensure docker is running for 'batch changes'
            try {
                await ensureDocker()
            } catch (error) {
                console.log(error)
                console.log('Docker required for batch changes')
                process.exit(1)
            }
            // set up batch change config
            const batchChange = batchChanges.releaseTrackingBatchChange(
                release.version,
                await batchChanges.sourcegraphCLIConfig()
            )

            // default values
            const notPatchRelease = release.patch === 0
            const versionRegex = '[0-9]+\\.[0-9]+\\.[0-9]+'
            const batchChangeURL = batchChanges.batchChangeURL(batchChange)
            const trackingIssue = await getTrackingIssue(await getAuthenticatedGitHubClient(), release)
            if (!trackingIssue) {
                throw new Error(`Tracking issue for version ${release.version} not found - has it been created yet?`)
            }

            // default PR content
            const defaultPRMessage = `release: sourcegraph@${release.version}`
            const prBodyAndDraftState = (
                actionItems: string[],
                customMessage?: string
            ): { draft: boolean; body: string } => {
                const defaultBody = `This pull request is part of the Sourcegraph ${release.version} release.
${customMessage || ''}

* [Release batch change](${batchChangeURL})
* ${trackingIssue ? `[Tracking issue](${trackingIssue.url})` : 'No tracking issue exists for this release'}`
                if (!actionItems || actionItems.length === 0) {
                    return { draft: false, body: defaultBody }
                }
                return {
                    draft: true, // further actions required before merge
                    body: `${defaultBody}

### :warning: Additional changes required

These steps must be completed before this PR can be merged, unless otherwise stated. Push any required changes directly to this PR branch.

${actionItems.map(item => `- [ ] ${item}`).join('\n')}

cc @${config.captainGitHubUsername}

### Test plan

CI checks in this repository should pass, and a manual review should confirm if the generated changes are correct.
`,
                }
            }

            const [previousVersion, nextVersion] = [
                `${previous.major}.${previous.minor}`,
                `${release.major}.${release.minor}`,
            ]

            // we join with escaped newlines for the 'sed' command
            const upgradeGuideEntry = [
                `## ${previousVersion} -> ${nextVersion}`,
                'TODO',
                hubSpotFeedbackFormStub(previousVersion),
            ].join('\\n\\n')

            // Render changes
            const createdChanges = await createChangesets({
                requiredCommands: ['comby', sed, 'find', 'go', 'src', 'sg'],
                changes: [
                    {
                        owner: 'sourcegraph',
                        repo: 'sourcegraph',
                        base: 'main',
                        head: `publish-${release.version}`,
                        commitMessage: notPatchRelease
                            ? `draft sourcegraph@${release.version} release`
                            : defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            // Update references to Sourcegraph versions in docs
                            `${sed} -i -E 's/version \`${versionRegex}\`/version \`${release.version}\`/g' doc/index.md`,
                            // Update sourcegraph/server:VERSION everywhere except changelog
                            `find . -type f -name '*.md' ! -name 'CHANGELOG.md' -exec ${sed} -i -E 's/sourcegraph\\/server:${versionRegex}/sourcegraph\\/server:${release.version}/g' {} +`,
                            // Update Sourcegraph versions in installation guides
                            `find ./doc/admin/deploy/ -type f -name '*.md' -exec ${sed} -i -E 's/SOURCEGRAPH_VERSION="v${versionRegex}"/SOURCEGRAPH_VERSION="v${release.version}"/g' {} +`,
                            `find ./doc/admin/deploy/ -type f -name '*.md' -exec ${sed} -i -E 's/--version ${versionRegex}/--version ${release.version}/g' {} +`,
                            // Update fork variables in installation guides
                            `find ./doc/admin/deploy/ -type f -name '*.md' -exec ${sed} -i -E "s/DEPLOY_SOURCEGRAPH_DOCKER_FORK_REVISION='v${versionRegex}'/DEPLOY_SOURCEGRAPH_DOCKER_FORK_REVISION='v${release.version}'/g" {} +`,
                            // Update sourcegraph.com frontpage
                            `${sed} -i -E 's/sourcegraph\\/server:${versionRegex}/sourcegraph\\/server:${release.version}/g' 'client/web/src/search/home/SelfHostInstructions.tsx'`,

                            notPatchRelease
                                ? `comby -in-place '{{$previousReleaseRevspec := ":[1]"}} {{$previousReleaseVersion := ":[2]"}} {{$currentReleaseRevspec := ":[3]"}} {{$currentReleaseVersion := ":[4]"}}' '{{$previousReleaseRevspec := ":[3]"}} {{$previousReleaseVersion := ":[4]"}} {{$currentReleaseRevspec := "v${release.version}"}} {{$currentReleaseVersion := "${release.major}.${release.minor}"}}' doc/_resources/templates/document.html`
                                : `comby -in-place 'currentReleaseRevspec := ":[1]"' 'currentReleaseRevspec := "v${release.version}"' doc/_resources/templates/document.html`,

                            // Update references to Sourcegraph deployment versions
                            `comby -in-place 'latestReleaseKubernetesBuild = newBuild(":[1]")' "latestReleaseKubernetesBuild = newBuild(\\"${release.version}\\")" cmd/frontend/internal/app/updatecheck/handler.go`,
                            `comby -in-place 'latestReleaseDockerServerImageBuild = newBuild(":[1]")' "latestReleaseDockerServerImageBuild = newBuild(\\"${release.version}\\")" cmd/frontend/internal/app/updatecheck/handler.go`,
                            `comby -in-place 'latestReleaseDockerComposeOrPureDocker = newBuild(":[1]")' "latestReleaseDockerComposeOrPureDocker = newBuild(\\"${release.version}\\")" cmd/frontend/internal/app/updatecheck/handler.go`,

                            // Support current release as the "previous release" going forward
                            notPatchRelease
                                ? `comby -in-place 'const minimumUpgradeableVersion = ":[1]"' 'const minimumUpgradeableVersion = "${release.version}"' enterprise/dev/ci/internal/ci/*.go`
                                : 'echo "Skipping minimumUpgradeableVersion bump on patch release"',

                            // Add a stub to add upgrade guide entries
                            notPatchRelease
                                ? `${sed} -i -E '/GENERATE UPGRADE GUIDE ON RELEASE/a \\\n\\n${upgradeGuideEntry}' doc/admin/updates/*.md`
                                : 'echo "Skipping upgrade guide entries on patch release"',
                        ],
                        ...prBodyAndDraftState(
                            ((): string[] => {
                                const items: string[] = []
                                if (notPatchRelease) {
                                    items.push('Update the upgrade guides in `doc/admin/updates`')
                                } else {
                                    items.push(
                                        'Update the [CHANGELOG](https://github.com/sourcegraph/sourcegraph/blob/main/CHANGELOG.md) to include all the changes included in this patch. Learn more about [how to update CHANGELOG.md](https://handbook.sourcegraph.com/departments/product-engineering/engineering/process/releases#changelogmd).',
                                        'If any specific upgrade steps are required, update the upgrade guides in `doc/admin/updates`'
                                    )
                                }
                                items.push(
                                    'Ensure all other pull requests in the batch change have been merged',
                                    'Run `yarn run release release:finalize` to generate the tags required. CI will not pass until this command is run.',
                                    'Re-run the build on this branch (using either `sg ci build --wait` or the Buildkite UI) and merge when the build passes.'
                                )
                                return items
                            })()
                        ),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'about',
                        base: 'main',
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            // Update sourcegraph/server:VERSION in all tsx files
                            `find . -type f -name '*.tsx' -exec ${sed} -i -E 's/sourcegraph\\/server:${versionRegex}/sourcegraph\\/server:${release.version}/g' {} +`,
                        ],
                        ...prBodyAndDraftState(
                            [],
                            notPatchRelease ? 'Note that this PR does *not* include the release blog post.' : undefined
                        ),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph',
                        base: `${release.major}.${release.minor}`,
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [`tools/update-docker-tags.sh ${release.version}`],
                        ...prBodyAndDraftState([]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-docker',
                        base: `${release.major}.${release.minor}`,
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [`tools/update-docker-tags.sh ${release.version}`],
                        ...prBodyAndDraftState([
                            `Follow the [release guide](https://github.com/sourcegraph/deploy-sourcegraph-docker/blob/master/RELEASING.md) to complete this PR ${
                                notPatchRelease ? '' : '(note: `pure-docker` release is optional for patch releases)'
                            }`,
                        ]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-aws',
                        base: 'master',
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            `${sed} -i -E 's/export SOURCEGRAPH_VERSION=${versionRegex}/export SOURCEGRAPH_VERSION=${release.version}/g' resources/amazon-linux2.sh`,
                        ],
                        ...prBodyAndDraftState([]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-digitalocean',
                        base: 'master',
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            `${sed} -i -E 's/export SOURCEGRAPH_VERSION=${versionRegex}/export SOURCEGRAPH_VERSION=${release.version}/g' resources/user-data.sh`,
                        ],
                        ...prBodyAndDraftState([]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-helm',
                        base: `release/${release.major}.${release.minor}`,
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            `for i in charts/*; do sg ops update-images -kind helm -pin-tag ${release.version} $i/.; done`,
                            `${sed} -i 's/appVersion:.*/appVersion: "${release.version}"/g' charts/*/Chart.yaml`,
                            `${sed} -i 's/version:.*/version: "${release.version}"/g' charts/*/Chart.yaml`,
                            './scripts/helm-docs.sh',
                        ],
                        ...prBodyAndDraftState([]),
                    },
                ],
                dryRun: dryRun.changesets,
            })

            // if changesets were actually published, set up a batch change and post in Slack
            if (!dryRun.changesets) {
                // Create batch change to track changes
                try {
                    console.log(`Creating batch change in ${batchChange.cliConfig.SRC_ENDPOINT}`)
                    await batchChanges.createBatchChange(
                        createdChanges,
                        batchChange,
                        `Track publishing of sourcegraph v${release.version}: ${trackingIssue?.url}`
                    )
                } catch (error) {
                    console.error(error)
                    console.error('Failed to create batch change for this release, continuing with announcement')
                }

                // Announce release update in Slack
                if (!dryRun.slack) {
                    await postMessage(
                        `:captain: *Sourcegraph ${release.version} has been staged.*

Batch change: ${batchChangeURL}`,
                        slackAnnounceChannel
                    )
                }
            }
        },
    },
    {
        id: 'release:add-to-batch-change',
        description: 'Manually add a change to a release batch change',
        argNames: ['changeRepo', 'changeID'],
        // Example: yarn run release release:add-to-batch-change sourcegraph/about 1797
        run: async (config, changeRepo, changeID) => {
            const { upcoming: release } = await releaseVersions(config)
            if (!changeRepo || !changeID) {
                throw new Error('Missing parameters (required: version, repo, change ID)')
            }

            const batchChange = batchChanges.releaseTrackingBatchChange(
                release.version,
                await batchChanges.sourcegraphCLIConfig()
            )
            await batchChanges.addToBatchChange(
                [
                    {
                        repository: changeRepo,
                        pullRequestNumber: parseInt(changeID, 10),
                    },
                ],
                batchChange
            )
            console.log(`Added ${changeRepo}#${changeID} to batch change ${batchChanges.batchChangeURL(batchChange)}`)
        },
    },
    {
        id: 'release:finalize',
        description: 'Run final tasks for sourcegraph/sourcegraph release pull requests',
        run: async config => {
            const { upcoming: release } = await releaseVersions(config)
            let failed = false

            // Push final tags
            const branch = `${release.major}.${release.minor}`
            const tag = `v${release.version}`
            for (const repo of ['deploy-sourcegraph', 'deploy-sourcegraph-docker']) {
                try {
                    await createTag(
                        await getAuthenticatedGitHubClient(),
                        {
                            owner: 'sourcegraph',
                            repo,
                            branch,
                            tag,
                        },
                        config.dryRun.tags || false
                    )
                } catch (error) {
                    console.error(error)
                    console.error(`Failed to create tag ${tag} on ${repo}@${branch}`)
                    failed = true
                }
            }

            if (failed) {
                throw new Error('Error occured applying some changes - please check log output')
            }
        },
    },
    {
        id: 'release:close',
        description: 'Mark a release as closed',
        run: async config => {
            const { slackAnnounceChannel, dryRun } = config
            const { upcoming: release } = await releaseVersions(config)
            const githubClient = await getAuthenticatedGitHubClient()

            // Create final GitHub release
            let githubRelease = ''
            try {
                githubRelease = await createLatestRelease(
                    githubClient,
                    {
                        owner: 'sourcegraph',
                        repo: 'sourcegraph',
                        release,
                    },
                    dryRun.tags
                )
            } catch (error) {
                console.error('Failed to generate GitHub release:', error)
                // Do not block process
            }

            // Set up announcement message
            const batchChangeURL = batchChanges.batchChangeURL(
                batchChanges.releaseTrackingBatchChange(release.version, await batchChanges.sourcegraphCLIConfig())
            )
            const releaseMessage = `*Sourcegraph ${release.version} has been published*

* Changelog: ${changelogURL(release.format())}
* GitHub release: ${githubRelease || 'No release generated'}
* Release batch change: ${batchChangeURL}`

            // Slack
            const slackMessage = `:captain: ${releaseMessage}`
            if (!dryRun.slack) {
                await postMessage(slackMessage, slackAnnounceChannel)
                console.log(`Posted to Slack channel ${slackAnnounceChannel}`)
            } else {
                console.log(`dryRun enabled, skipping Slack post to ${slackAnnounceChannel}: ${slackMessage}`)
            }

            // GitHub tracking issues
            const trackingIssue = await getTrackingIssue(githubClient, release)
            if (!trackingIssue) {
                console.warn(`Could not find tracking issue for release ${release.version} - skipping`)
            } else {
                // Note patch release requests if there are any outstanding
                let comment = `${releaseMessage}

@${config.captainGitHubUsername}: Please complete the post-release steps before closing this issue.`
                const patchRequestIssues = await queryIssues(githubClient, '', [IssueLabel.PATCH_REQUEST])
                if (patchRequestIssues && patchRequestIssues.length > 0) {
                    comment += `
Please also update outstanding patch requests, if relevant:

${patchRequestIssues.map(issue => `* #${issue.number}`).join('\n')}`
                }
                if (!dryRun.trackingIssues) {
                    const commentURL = await commentOnIssue(githubClient, trackingIssue, comment)
                    console.log(`Please make sure to follow up on the release issue: ${commentURL}`)
                } else {
                    console.log(`dryRun enabled, skipping GitHub comment to ${trackingIssue.url}: ${comment}`)
                }
            }
        },
    },
    {
        id: 'util:clear-cache',
        description: 'Clear release tool cache',
        run: () => {
            rmdirSync(cacheFolder, { recursive: true })
        },
    },
    {
        id: '_test:google-calendar',
        description: 'Test Google Calendar integration',
        run: async config => {
            const googleCalendar = await getClient()
            await ensureEvent(
                {
                    title: 'TEST EVENT',
                    startDateTime: new Date(config.releaseDate).toISOString(),
                    endDateTime: addMinutes(new Date(config.releaseDate), 1).toISOString(),
                    transparency: 'transparent',
                },
                googleCalendar
            )
        },
    },
    {
        id: '_test:slack',
        description: 'Test Slack integration',
        argNames: ['channel', 'message'],
        run: async ({ dryRun }, channel, message) => {
            if (!dryRun.slack) {
                await postMessage(message, channel)
            }
        },
    },
    {
        id: '_test:batchchange-create-from-changes',
        description: 'Test batch changes integration',
        argNames: ['batchchangeConfigJSON'],
        // Example: yarn run release _test:batchchange-create-from-changes "$(cat ./.secrets/test-batch-change-import.json)"
        run: async (_config, batchchangeConfigJSON) => {
            const batchChangeConfig = JSON.parse(batchchangeConfigJSON) as {
                changes: CreatedChangeset[]
                name: string
                description: string
            }

            // set up src-cli
            await commandExists('src')
            const batchChange = {
                name: batchChangeConfig.name,
                description: batchChangeConfig.description,
                namespace: 'sourcegraph',
                cliConfig: await batchChanges.sourcegraphCLIConfig(),
            }

            await batchChanges.createBatchChange(
                batchChangeConfig.changes,
                batchChange,
                'release tool testing batch change'
            )
            console.log(`Created batch change ${batchChanges.batchChangeURL(batchChange)}`)
        },
    },
    {
        id: '_test:config',
        description: 'Test release configuration loading',
        run: config => {
            console.log(JSON.stringify(config, null, '  '))
        },
    },
    {
        id: '_test:dockerensure',
        description: 'test docker ensure function',
        run: async () => {
            try {
                await ensureDocker()
            } catch (error) {
                console.log(error)
                console.log('Docker required for batch changes')
                process.exit(1)
            }
        },
    },
]
