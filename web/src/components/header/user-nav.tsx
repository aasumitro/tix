import { LogOut, Settings} from "lucide-react"

import { Avatar, AvatarFallback } from "../ui/avatar"
import { Button } from "../ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import {useState} from 'react';
import {LogoutModal} from '../logout-modal';
import {User} from '../../libs/model/user';
import {SettingModal} from '../setting-modal';

interface UserNavProps {
  user?: User;
}

export function UserNav(props: UserNavProps) {
  const [showProfileDialog, setShowProfileDialog] = useState(false);
  const [showLogoutDialog, setShowLogoutDialog] = useState(false);

  const logoutCallback = () => setShowLogoutDialog(false);

  const profileCallback = () => setShowProfileDialog(false);

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="relative h-8 w-8 rounded-full">
            <Avatar className="h-8 w-8">
              <AvatarFallback>{props?.user?.username?.substring(0, 2) ?? "-"}</AvatarFallback>
            </Avatar>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56" align="end" forceMount>
          <DropdownMenuLabel className="font-normal">
            <div className="flex flex-col space-y-1">
              <p className="text-sm font-medium leading-none">
                @{props?.user?.username ?? "-"}
              </p>
              <p className="text-xs leading-none text-muted-foreground">
                {props?.user?.email ?? "-"}
              </p>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem onClick={() => setShowProfileDialog(true)}>
              <Settings className="mr-2 h-4 w-4" />
              <span>Settings</span>
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={() => setShowLogoutDialog(true)}>
            <LogOut className="mr-2 h-4 w-4" />
            <span>Log out</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <SettingModal
        show={showProfileDialog}
        callback={profileCallback} />

      <LogoutModal
        show={showLogoutDialog}
        callback={logoutCallback} />
    </>
  )
}