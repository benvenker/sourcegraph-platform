import { storiesOf } from '@storybook/react'

import { getDocumentNode } from '@sourcegraph/http-client'
import { MockedTestProvider } from '@sourcegraph/shared/src/testing/apollo'

import { WebStory } from '../../../components/WebStory'
import {
    BatchChangesCredentialFields,
    CheckBatchChangesCredentialResult,
    ExternalServiceKind,
} from '../../../graphql-operations'

import { CHECK_BATCH_CHANGES_CREDENTIAL } from './backend'
import { CodeHostConnectionNode } from './CodeHostConnectionNode'

const { add } = storiesOf('web/batches/settings/CodeHostConnectionNode', module).addDecorator(story => (
    <div className="p-3 container">{story()}</div>
))

const checkCredResult = (): CheckBatchChangesCredentialResult => ({
    checkBatchChangesCredential: {
        alwaysNil: null,
    },
})

const sshCredential = (isSiteCredential: boolean): BatchChangesCredentialFields => ({
    id: '123',
    isSiteCredential,
    sshPublicKey:
        'rsa-ssh randorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorandorando',
})

add('Overview', () => (
    <WebStory>
        {props => (
            <MockedTestProvider
                mocks={[
                    {
                        request: {
                            query: getDocumentNode(CHECK_BATCH_CHANGES_CREDENTIAL),
                            variables: {
                                id: '123',
                            },
                        },
                        result: {
                            data: checkCredResult(),
                        },
                        // Some sort of delay to see the spinner
                        delay: 1000,
                    },
                ]}
            >
                <CodeHostConnectionNode
                    {...props}
                    node={{
                        credential: sshCredential(false),
                        externalServiceKind: ExternalServiceKind.GITHUB,
                        externalServiceURL: 'https://github.com/',
                        requiresSSH: false,
                        requiresUsername: false,
                    }}
                    refetchAll={() => {}}
                    userID="123"
                />
            </MockedTestProvider>
        )}
    </WebStory>
))
