import React from 'react'

import classNames from 'classnames'
import PlusIcon from 'mdi-react/PlusIcon'

import { Position, Menu, MenuButton, MenuList, MenuLink, Icon, Link } from '@sourcegraph/wildcard'

import { AuthenticatedUser } from '../../auth'

import { CreateAction } from './createActions'

import createActionsStyles from './CreateActions.module.scss'
import navStyles from './SearchResultsInfoBar.module.scss'

export interface CreateActionsMenuProps {
    createActions: CreateAction[]
    createCodeMonitorAction: CreateAction | null
    canCreateMonitor: boolean
    authenticatedUser: Pick<AuthenticatedUser, 'id'> | null
}

export const CreateActionsMenu: React.FunctionComponent<CreateActionsMenuProps> = ({
    createActions,
    createCodeMonitorAction,
    canCreateMonitor,
    authenticatedUser,
}) => (
    <Menu>
        {({ isExpanded }) => (
            <>
                <li className={classNames('mr-2', createActionsStyles.menu, navStyles.navItem)}>
                    <MenuButton
                        className={classNames('d-flex align-items-center text-decoration-none')}
                        aria-label={`${isExpanded ? 'Close' : 'Open'} create actions menu`}
                        variant="secondary"
                        outline={true}
                        size="sm"
                    >
                        <Icon role="img" aria-hidden={true} className="mr-1" as={PlusIcon} />
                        Create …
                    </MenuButton>
                </li>
                <MenuList position={Position.bottomStart} aria-label="Create Actions. Open menu">
                    {createActions.map(createAction => (
                        <MenuLink key={createAction.label} as={Link} to={createAction.url}>
                            <Icon role="img" aria-hidden="true" className="mr-1" as={createAction.icon} />
                            {createAction.label}
                        </MenuLink>
                    ))}
                    {createCodeMonitorAction && (
                        <MenuLink
                            as={Link}
                            disabled={!authenticatedUser || !canCreateMonitor}
                            data-tooltip={
                                authenticatedUser && !canCreateMonitor
                                    ? 'Code monitors only support type:diff or type:commit searches.'
                                    : undefined
                            }
                            to={createCodeMonitorAction.url}
                        >
                            <Icon role="img" aria-hidden={true} className="mr-1" as={createCodeMonitorAction.icon} />
                            Create Monitor
                        </MenuLink>
                    )}
                </MenuList>
            </>
        )}
    </Menu>
)
