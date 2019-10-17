import { LoadingSpinner } from '@sourcegraph/react-loading-spinner'
import React, { Suspense } from 'react'
import { Redirect, Route, RouteComponentProps, Switch } from 'react-router'
import { Observable } from 'rxjs'
import { ActivationProps } from '../../shared/src/components/activation/Activation'
import { FetchFileCtx } from '../../shared/src/components/CodeExcerpt'
import { ExtensionsControllerProps } from '../../shared/src/extensions/controller'
import * as GQL from '../../shared/src/graphql/schema'
import { ResizablePanel } from '../../shared/src/panel/Panel'
import { PlatformContextProps } from '../../shared/src/platform/context'
import { SettingsCascadeProps } from '../../shared/src/settings/settings'
import { ErrorLike } from '../../shared/src/util/errors'
import { parseHash } from '../../shared/src/util/url'
import { ErrorBoundary } from './components/ErrorBoundary'
import { useScrollToLocationHash } from './components/useScrollToLocationHash'
import { GlobalContributions } from './contributions'
import { ExploreSectionDescriptor } from './explore/ExploreArea'
import { ExtensionAreaRoute } from './extensions/extension/ExtensionArea'
import { ExtensionAreaHeaderNavItem } from './extensions/extension/ExtensionAreaHeader'
import { ExtensionsAreaRoute } from './extensions/ExtensionsArea'
import { ExtensionsAreaHeaderActionButton } from './extensions/ExtensionsAreaHeader'
import { GlobalAlerts } from './global/GlobalAlerts'
import { GlobalDebug } from './global/GlobalDebug'
import { KEYBOARD_SHORTCUT_SHOW_HELP, KeyboardShortcutsProps } from './keyboardShortcuts/keyboardShortcuts'
import { KeyboardShortcutsHelp } from './keyboardShortcuts/KeyboardShortcutsHelp'
import { IntegrationsToast } from './marketing/IntegrationsToast'
import { GlobalNavbar } from './nav/GlobalNavbar'
import { OrgAreaRoute } from './org/area/OrgArea'
import { OrgAreaHeaderNavItem } from './org/area/OrgHeader'
import { fetchHighlightedFileLines } from './repo/backend'
import { RepoContainerRoute } from './repo/RepoContainer'
import { RepoHeaderActionButton } from './repo/RepoHeader'
import { RepoRevContainerRoute } from './repo/RepoRevContainer'
import { LayoutRouteProps } from './routes'
import { parseSearchURLQuery, PatternTypeProps } from './search'
import { SiteAdminAreaRoute } from './site-admin/SiteAdminArea'
import { SiteAdminSideBarGroups } from './site-admin/SiteAdminSidebar'
import { ThemePreferenceProps, ThemeProps } from './theme'
import { EventLogger, EventLoggerProps } from './tracking/eventLogger'
import { UserAreaRoute } from './user/area/UserArea'
import { UserAreaHeaderNavItem } from './user/area/UserAreaHeader'
import { UserSettingsAreaRoute } from './user/settings/UserSettingsArea'
import { UserSettingsSidebarItems } from './user/settings/UserSettingsSidebar'
import { parseBrowserRepoURL } from './util/url'
import LiteralSearchToast from './marketing/LiteralSearchToast'

export interface LayoutProps
    extends RouteComponentProps<any>,
        SettingsCascadeProps,
        PlatformContextProps,
        ExtensionsControllerProps,
        KeyboardShortcutsProps,
        ThemeProps,
        EventLoggerProps,
        ThemePreferenceProps,
        ActivationProps,
        PatternTypeProps {
    exploreSections: readonly ExploreSectionDescriptor[]
    extensionAreaRoutes: readonly ExtensionAreaRoute[]
    extensionAreaHeaderNavItems: readonly ExtensionAreaHeaderNavItem[]
    extensionsAreaRoutes: readonly ExtensionsAreaRoute[]
    extensionsAreaHeaderActionButtons: readonly ExtensionsAreaHeaderActionButton[]
    siteAdminAreaRoutes: readonly SiteAdminAreaRoute[]
    siteAdminSideBarGroups: SiteAdminSideBarGroups
    siteAdminOverviewComponents: readonly React.ComponentType[]
    userAreaHeaderNavItems: readonly UserAreaHeaderNavItem[]
    userAreaRoutes: readonly UserAreaRoute[]
    userSettingsSideBarItems: UserSettingsSidebarItems
    userSettingsAreaRoutes: readonly UserSettingsAreaRoute[]
    orgAreaHeaderNavItems: readonly OrgAreaHeaderNavItem[]
    orgAreaRoutes: readonly OrgAreaRoute[]
    repoContainerRoutes: readonly RepoContainerRoute[]
    repoRevContainerRoutes: readonly RepoRevContainerRoute[]
    repoHeaderActionButtons: readonly RepoHeaderActionButton[]
    routes: readonly LayoutRouteProps[]

    authenticatedUser: GQL.IUser | null

    /**
     * The subject GraphQL node ID of the viewer, which is used to look up the viewer's settings. This is either
     * the site's GraphQL node ID (for anonymous users) or the authenticated user's GraphQL node ID.
     */
    viewerSubject: Pick<GQL.ISettingsSubject, 'id' | 'viewerCanAdminister'>

    telemetryService: EventLogger

    // Search
    navbarSearchQuery: string
    onNavbarQueryChange: (query: string) => void
    fetchHighlightedFileLines: (ctx: FetchFileCtx, force?: boolean) => Observable<string[]>
    searchRequest: (
        query: string,
        version: string,
        patternType: GQL.SearchPatternType,
        { extensionsController }: ExtensionsControllerProps<'services'>
    ) => Observable<GQL.ISearchResults | ErrorLike>
    isSourcegraphDotCom: boolean
    showCampaigns: boolean
    children?: never
}

export const Layout: React.FunctionComponent<LayoutProps> = props => {
    const isSearchHomepage = props.location.pathname === '/search' && !parseSearchURLQuery(props.location.search)

    const needsSiteInit = window.context.showOnboarding
    const isSiteInit = props.location.pathname === '/site-admin/init'

    const showLiteralSearchToast = props.isSourcegraphDotCom
        ? true
        : props.authenticatedUser &&
          props.authenticatedUser.createdAt &&
          userWasCreatedBeforeSourcegraph39Release(props.authenticatedUser.createdAt)

    useScrollToLocationHash(props.location)
    // Remove trailing slash (which is never valid in any of our URLs).
    if (props.location.pathname !== '/' && props.location.pathname.endsWith('/')) {
        return <Redirect to={{ ...props.location, pathname: props.location.pathname.slice(0, -1) }} />
    }

    return (
        <div className="layout">
            <KeyboardShortcutsHelp
                keyboardShortcutForShow={KEYBOARD_SHORTCUT_SHOW_HELP}
                keyboardShortcuts={props.keyboardShortcuts}
            />
            <GlobalAlerts
                isSiteAdmin={!!props.authenticatedUser && props.authenticatedUser.siteAdmin}
                settingsCascade={props.settingsCascade}
            />
            {!needsSiteInit && !isSiteInit && !!props.authenticatedUser && (
                <IntegrationsToast history={props.history} />
            )}
            {!isSiteInit && showLiteralSearchToast && (
                <LiteralSearchToast isSourcegraphDotCom={props.isSourcegraphDotCom} />
            )}
            {!isSiteInit && <GlobalNavbar {...props} lowProfile={isSearchHomepage} />}
            {needsSiteInit && !isSiteInit && <Redirect to="/site-admin/init" />}
            <ErrorBoundary location={props.location}>
                <Suspense fallback={<LoadingSpinner className="icon-inline m-2" />}>
                    <Switch>
                        {/* eslint-disable react/jsx-no-bind */}
                        {props.routes.map(({ render, condition = () => true, ...route }) => {
                            const isFullWidth = !route.forceNarrowWidth
                            return (
                                condition(props) && (
                                    <Route
                                        {...route}
                                        key="hardcoded-key" // see https://github.com/ReactTraining/react-router/issues/4578#issuecomment-334489490
                                        component={undefined}
                                        render={routeComponentProps => (
                                            <div
                                                className={[
                                                    'layout__app-router-container',
                                                    `layout__app-router-container--${
                                                        isFullWidth ? 'full-width' : 'restricted'
                                                    }`,
                                                ].join(' ')}
                                            >
                                                {render({ ...props, ...routeComponentProps })}
                                            </div>
                                        )}
                                    />
                                )
                            )
                        })}
                        {/* eslint-enable react/jsx-no-bind */}
                    </Switch>
                </Suspense>
            </ErrorBoundary>
            {parseHash(props.location.hash).viewState && props.location.pathname !== '/sign-in' && (
                <ResizablePanel
                    {...props}
                    repoName={`git://${parseBrowserRepoURL(props.location.pathname).repoName}`}
                    fetchHighlightedFileLines={fetchHighlightedFileLines}
                />
            )}
            <GlobalContributions
                key={3}
                extensionsController={props.extensionsController}
                platformContext={props.platformContext}
                history={props.history}
            />
            <GlobalDebug {...props} />
        </div>
    )
}

// Whether an authenticated user account was created before our 3.9 release date.
// Used to determine whether we display the literal search toast. We don't want to show
// it to users created after the release date.
function userWasCreatedBeforeSourcegraph39Release(createdAtDate: string): boolean {
    // One day before the release.
    const SOURCEGRAPH_39_RELEASE = Date.parse('2019-10-19')
    const parsedCreatedAt = Date.parse(createdAtDate)

    return parsedCreatedAt < SOURCEGRAPH_39_RELEASE
}
