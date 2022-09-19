import React, { ReactElement } from 'react'
import { buttonDisplayInDesignSystem } from '../view/button'

/* Pages / View Controllers */

// View Controllers should be the mediator between interactor or router and views
// When user action through view or other events request or change data,
//  it should call the relevant interactor function.
// When user action through view or other events cause page to change,
//  it should call the relevant router function.

export const DesignSystemPage = (): ReactElement => {
  // Use this to experiment with UI
  return (
    <div className='bg-neutral-a max-w-lg m-auto p-2'>
      {buttonDisplayInDesignSystem()}
    </div>
  )
}
