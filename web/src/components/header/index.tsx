import EventSwitcher from './event-switcher';
import {MainNav} from './main-nav';
import {Search} from './search';
import {UserNav} from './user-nav';
import * as React from 'react';
import {User} from '../../libs/model/user';

interface HeaderProps {
  user?: User
}

export function Header(props: HeaderProps) {
  return (
    <>
      <header>
        <div className="border-b">
          <div className="flex h-16 items-center px-4">
            <EventSwitcher />
            <MainNav className="mx-6" />
            <div className="ml-auto flex items-center space-x-4">
              <Search />
              <UserNav user={props.user}/>
            </div>
          </div>
        </div>
      </header>
    </>
  )
}