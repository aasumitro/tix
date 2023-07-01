import * as React from 'react';
import {useEffect, useState} from 'react';
import {BaseUrl, Endpoint} from '../libs/api';
import {User} from '../libs/model/user';
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle
} from './ui/alert-dialog';
import { Moon, Sun, XIcon} from 'lucide-react';
import {Label} from './ui/label';
import {Input} from './ui/input';
import {Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue} from './ui/select';

interface SettingModalProps {
  show: boolean,
  callback(): void,
}

export function SettingModal(props: SettingModalProps) {
  const [showSettingDialog, setShowSettingDialog] = React.useState(false)
  const [user, setUser] = useState<User>();

  useEffect(() => {
    setShowSettingDialog(props.show)
    getUserProfile()
  }, [props.show])

  function getUserProfile() {
    fetch(`${BaseUrl}/${Endpoint.Auth.Profile}`, {
      method: 'GET',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include'
    }).then((resp) => {
      if (resp.status === 401) {
        localStorage.removeItem("current_event")
        localStorage.removeItem("is_login")
        window.location.href = "/admin"
        return
      }
      return resp.json()
    }).then((resp) => setUser(resp.data))
  }

  return (<AlertDialog
    open={showSettingDialog}
    onOpenChange={() => {
      setShowSettingDialog(false)
      props.callback()
    }}
  >
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>
          Settings
        </AlertDialogTitle>
        <AlertDialogDescription>
          Your personalized hub for managing your online experience! have complete control over your profile and preferences.
        </AlertDialogDescription>
        <button className="absolute top-0 right-0 pr-4 pt-4" onClick={() => {
          setShowSettingDialog(false)
          props.callback()
        }}>
          <XIcon/>
        </button>
      </AlertDialogHeader>

      <div className="space-y-4 py-2 pb-4 text-left">
        <div className="space-y-2">
          <Label htmlFor="uuid">UUID</Label>
          <Input
            id="uuid"
            placeholder="admin@tix.id"
            value={user?.uuid}
            disabled/>
        </div>
        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            placeholder="admin@tix.id"
            value={user?.email}
            disabled/>
        </div>
        <div className="space-y-2">
          <Label htmlFor="username">Username</Label>
          <Input
            id="username"
            placeholder="@tix_id"
            value={`@${user?.username}`}
            disabled/>
        </div>
        <div className="flex flex-row gap-2 justify-between">
          <div className="space-y-2 w-full">
            <Label htmlFor="username">Theme</Label>
            <Select disabled value="light">
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select a theme" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Theme</SelectLabel>
                  <SelectItem value="light">
                    <div className="flex flex-row gap-2 items-center">
                      <Sun className="w-4 h-4"/>
                      Light
                    </div>
                  </SelectItem>
                  <SelectItem value="dark">
                    <div className="flex flex-row gap-2 items-center">
                      <Moon className="w-4 h-4"/>
                      Light
                    </div>
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2 w-full">
            <Label htmlFor="language">Language</Label>
            <Select disabled value="en">
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select a language" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Language</SelectLabel>
                  <SelectItem value="id">
                    <div className="flex flex-row gap-2 items-center">
                      <svg className="w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="40" zoomAndPan="magnify" viewBox="0 0 30 30.000001" height="40" preserveAspectRatio="xMidYMid meet" version="1.0"><defs><clipPath id="id1"><path d="M 2.128906 5.222656 L 27.53125 5.222656 L 27.53125 15 L 2.128906 15 Z M 2.128906 5.222656 " clipRule="nonzero"/></clipPath><clipPath id="id2"><path d="M 2.128906 14 L 27.53125 14 L 27.53125 23.371094 L 2.128906 23.371094 Z M 2.128906 14 " clipRule="nonzero"/></clipPath></defs><g clipPath="url(#id1)"><path fill="rgb(86.268616%, 12.159729%, 14.898682%)" d="M 24.703125 5.222656 L 4.957031 5.222656 C 3.398438 5.222656 2.132812 6.472656 2.132812 8.015625 L 2.132812 14.296875 L 27.523438 14.296875 L 27.523438 8.015625 C 27.523438 6.472656 26.261719 5.222656 24.703125 5.222656 Z M 24.703125 5.222656 " fillOpacity="1" fillRule="nonzero"/></g><g clipPath="url(#id2)"><path fill="rgb(93.328857%, 93.328857%, 93.328857%)" d="M 27.523438 20.578125 C 27.523438 22.121094 26.261719 23.371094 24.703125 23.371094 L 4.957031 23.371094 C 3.398438 23.371094 2.132812 22.121094 2.132812 20.578125 L 2.132812 14.296875 L 27.523438 14.296875 Z M 27.523438 20.578125 " fillOpacity="1" fillRule="nonzero"/></g></svg>
                      Indonesia
                    </div>
                  </SelectItem>
                  <SelectItem value="en">
                    <div className="flex flex-row gap-2 items-center">
                      <svg className="w-4 h-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 52 34">
                        <g fillRule="evenodd">
                          <path fill="#b22234" d="M0 0h52v34H0z"/>
                          <path fill="#3c3b6e" d="M0 0h19v13H0zm0 21h19v13H0zm33 0h19v13H33z"/>
                          <path fill="#fff" d="M0 13h52v8H0z"/>
                        </g>
                        <g stroke="#fff" strokeWidth="4">
                          <path d="M6 0v34M0 6h52M0 17h52M0 28h52"/>
                        </g>
                      </svg>
                      English
                    </div>
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>
    </AlertDialogContent>
  </AlertDialog>)
}