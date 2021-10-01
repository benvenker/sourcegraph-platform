import { Observable, of } from 'rxjs'
import { map } from 'rxjs/operators'

import { dataOrThrowErrors, gql } from '@sourcegraph/shared/src/graphql/graphql'
import { memoizeObservable } from '@sourcegraph/shared/src/util/memoizeObservable'
import { ParsedRepoURI, makeRepoURI } from '@sourcegraph/shared/src/util/url'

import { requestGraphQL } from '../../backend/graphql'
import { BlobFileFields, BlobResult, BlobVariables } from '../../graphql-operations'

function fetchBlobCacheKey(parsed: ParsedRepoURI & { disableTimeout: boolean }): string {
    return makeRepoURI(parsed) + String(parsed.disableTimeout)
}

export const fetchBlob = memoizeObservable(
    (args: {
        repoName: string
        commitID: string
        filePath: string
        disableTimeout: boolean
    }): Observable<BlobFileFields | null> => {
        const cacheKey = new URLSearchParams(window.location.search).get('cache')
        const cachedBlob = window.blob as string
        console.log('cacheKey', cacheKey)

        if (cachedBlob && cacheKey) {
            return of({
                content: '',
                richHTML: '',
                highlight: {
                    aborted: false,
                    html: cachedBlob,
                },
            })
        }

        return requestGraphQL<BlobResult, BlobVariables>(
            gql`
                query Blob($repoName: String!, $commitID: String!, $filePath: String!, $disableTimeout: Boolean!) {
                    repository(name: $repoName) {
                        commit(rev: $commitID) {
                            file(path: $filePath) {
                                ...BlobFileFields
                            }
                        }
                    }
                }

                fragment BlobFileFields on File2 {
                    content
                    richHTML
                    highlight(disableTimeout: $disableTimeout) {
                        aborted
                        html
                    }
                }
            `,
            args
        ).pipe(
            map(dataOrThrowErrors),
            map(data => {
                if (!data.repository?.commit) {
                    throw new Error('Commit not found')
                }
                return data.repository.commit.file
            })
        )
    },
    fetchBlobCacheKey
)
