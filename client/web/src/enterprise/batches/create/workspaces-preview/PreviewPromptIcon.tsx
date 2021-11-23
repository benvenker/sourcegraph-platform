import React, { memo, SVGProps } from 'react'

export const PreviewPromptIcon: React.FunctionComponent<SVGProps<SVGSVGElement>> = memo(props => (
    <svg width="50" height="50" viewBox="0 0 50 50" fill="none" xmlns="http://www.w3.org/2000/svg" {...props}>
        <circle cx="25" cy="25" r="24.5" fill="#EFF2F5" stroke="#E6EBF2" />
        <path
            d="M20.6452 21.9412V22.296H21H22.3333H22.6882V21.9412V19.4706V19.1158H22.3333H21H20.6452V19.4706V21.9412ZM15.6667 19.1158H15.3118V19.4706V21.9412V22.296H15.6667H18.3333H18.6882V21.9412V19.4706V19.1158H18.3333H15.6667ZM23.6667 33.9393H20.0215V32.1784H23.6667H24.0215V31.8235V28.4725H25.9785V31.8235V32.1784H26.3333H29.9785V33.9393H26.3333H25.9785V34.2941V37.6452H24.0215V34.2941V33.9393H23.6667ZM14.3333 17.3548H35.6667C35.9346 17.3548 36.1865 17.4537 36.3683 17.6221C36.5492 17.7897 36.6452 18.0109 36.6452 18.2353V23.1765C36.6452 23.4009 36.5492 23.6221 36.3683 23.7897C36.1865 23.9581 35.9346 24.0569 35.6667 24.0569H14.3333C14.0654 24.0569 13.8135 23.9581 13.6317 23.7897C13.4508 23.6221 13.3548 23.4009 13.3548 23.1765V18.2353C13.3548 18.0109 13.4508 17.7897 13.6317 17.6221C13.8135 17.4537 14.0654 17.3548 14.3333 17.3548Z"
            fill="#DBE2F0"
            stroke="#A6B6D9"
            strokeWidth="0.70965"
        />
    </svg>
))
