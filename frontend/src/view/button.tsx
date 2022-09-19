import React, { ReactElement } from 'react'

/* Views */

// Views should only contain presentation logic for themselves or their children.
//  Event should be handled by parent (e.g. view controller/page) passing down event handler to the view.
//  example (in view controller): createView(initialConfig, eventHandler)
// Views should not know anything about view controllers/pages
//  Don't pass in view controller into view like this: function createView(viewController){}

const designSystemDisplayList: ReactElement[] = []
export function buttonDisplayInDesignSystem (): ReactElement {
  return <>
    {designSystemDisplayList
      .map((view, index) => <div className='design-system-item' key={index}>{view}</div>)}
  </>
}

interface ButtonProps {
  children: string
  onClick?: () => void
}

export const RoundedButton = function (props: ButtonProps): ReactElement {
  return (
    <button
      className={'rounded-full px-3 py-1 bg-main-c text-neutral-f hover:bg-main-b active:bg-main-a'}
      onClick={props.onClick}
    >
      {props.children}
    </button>
  )
}
designSystemDisplayList.push(
  <RoundedButton
    onClick={() => { alert('RoundedButton') }}
  >Example RoundedButton</RoundedButton>
)

export const BorderlessButton = function (props: ButtonProps): ReactElement {
  return (
    <button
      className={'underline text-main-d hover:no-underline hover:text-main-e active:no-underline active:text-main-c'}
      onClick={props.onClick}
    >
      {props.children}
    </button>
  )
}
designSystemDisplayList.push(
  <BorderlessButton
    onClick={() => { alert('BorderlessButton') }}
  >Example BorderlessButton</BorderlessButton>
)
