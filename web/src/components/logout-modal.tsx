import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent, AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader, AlertDialogTitle,
} from './ui/alert-dialog';
import * as React from 'react';
import {useEffect} from 'react';
import {BaseUrl, Endpoint} from '../libs/api';

interface LogoutModalProps {
  show: boolean,
  callback(): void,
}

export function LogoutModal(props: LogoutModalProps) {
  const [showLogoutDialog, setShowLogoutDialog] = React.useState(false)

  useEffect(() => {
    setShowLogoutDialog(props.show)
  }, [props.show])

  async function doLogout(isLoggedOut: boolean) {
    if (isLoggedOut) {
      const request = await fetch(`${BaseUrl}/${Endpoint.Auth.Logout}`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
      });
      const response = await request
      if (response.status === 401) {
        localStorage.removeItem("current_event")
        localStorage.removeItem("is_login")
        window.location.href = "/admin"
      }
    }
    setShowLogoutDialog(false)
    props.callback()
  }

  return (
    <AlertDialog
      open={showLogoutDialog}
      onOpenChange={() => doLogout(false)}
    >
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            Logout
          </AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to logout?
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction onClick={() => doLogout(true)}>Lemme Out</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}