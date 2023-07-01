import {User} from '../libs/model/user';
import {Button} from './ui/button';
import {Filter, Loader,  Mail, Trash, Verified} from 'lucide-react';
import {Input} from './ui/input';
import {
  DropdownMenu, DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from './ui/dropdown-menu';
import {Card, CardContent} from './ui/card';
import {Avatar, AvatarFallback} from './ui/avatar';
import {Badge} from './ui/badge';
import * as React from 'react';
import {InviteNewUserModal} from './invite-new-user';
import {BaseUrl, Endpoint} from '../libs/api';
import {toast} from './ui/use-toast';
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel,
  AlertDialogContent, AlertDialogDescription, AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger
} from './ui/alert-dialog';
import {useEffect} from 'react';

interface UserListDataProps {
  users: User[];
  doRefreshCallback: () => void;
}

interface TableData {
  id: string;
  check: boolean;
}

const status: TableData[] = [
  {id: "verified", check: false},
  {id: "need confirmation", check: false},
]

export function UserListData(props: UserListDataProps) {
  const [showInviteUserDialog, setShowInviteUserDialog] = React.useState(false)
  const [userTemp, setUserTemp] = React.useState<User[]>([])

  useEffect(() => {setUserTemp(props.users)}, [props.users])

  function inviteMemberCallback() {
    setShowInviteUserDialog(false)
    props.doRefreshCallback()
  }

  function deleteAccount(uuid: string) {
    fetch(`${BaseUrl}/${Endpoint.User.Remove}/${uuid}`, {
      method: 'DELETE',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
    })
      .then((resp) => {
        if (resp.status === 204) {
          toast({
            variant: "default",
            title: "Account deleted!",
            description: "Account has been deleted successfully.",
          })
          props.doRefreshCallback()
          return
        }
        return resp.json();
      })
      .then(resp => {
        if (!resp) {
          return
        }

        if (resp.code === 406 || resp.code === 400) {
          toast({
            variant: "destructive",
            title: "Failed to delete account!",
            description: resp.data,
          })
          return
        }
      })
  }

  function filterUserByEmail(email: string) {
    setUserTemp((email !== "")
      ? props.users.filter((user) => user.email.includes(email))
      : props.users)
  }

  return (<>
    <div className="container w-full h-full space-y-4 py-12 flex flex-col">
      <div className="flex flex-row items-center justify-between mb-10">
        <div>
          <h1 className="text-2xl font-bold">Users</h1>
          <p>List of invited users who can access this site.</p>
        </div>
        <Button onClick={() => setShowInviteUserDialog(true)}>
          <Mail className="mr-2 h-4 w-4" /> Invite via email
        </Button>
      </div>

      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-1 items-center space-x-2">
          <Input
            placeholder="Search by email . . ."
            className="h-8 w-[150px] lg:w-[250px]"
            onChange={(e) => filterUserByEmail(e.currentTarget.value)}
          />
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="outline"
              size="sm"
              className="ml-auto h-8 flex"
            >
              <Filter className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-[250px]">
            <DropdownMenuLabel>Filter data</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {status.map((column) => {
              return (
                <DropdownMenuCheckboxItem
                  key={column.id}
                  className="capitalize"
                  checked={column.check}
                >
                  {column.id}
                </DropdownMenuCheckboxItem>
              )
            })}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <div className="flex flex-col pt-4 gap-4">
        {userTemp.map((user) => (
          <Card key={user.id} className="rounded-md">
            <CardContent className="flex flex-row justify-between items-center px-4 py-2">
              <div className="flex w-40">
                <Avatar className="h-8 w-8">
                  <AvatarFallback>{user.username.substring(0, 2)}</AvatarFallback>
                </Avatar>
                <div className="flex flex-row items-center gap-2 ml-4">
                  <h5 className="text-md font-bold text-gray-600">@{user.username}</h5>
                  <p className="text-sm text-gray-500">{user.email}</p>
                </div>
              </div>
              <div className="flex gap-2">
                {user.is_verified && <>
                    <Badge variant="destructive">
                        Admin
                    </Badge>
                    <Badge variant="default">
                        <Verified  className="h-4 w-4 text-green-300 mr-2"/>
                        Verified
                    </Badge>
                </>}
                {!user.is_verified && <Badge variant="secondary">
                    <Loader className="h-3 w-3 animate-spin mr-2" />
                    waiting email confirmation . . .
                </Badge>}
              </div>
              <div className="flex">
                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button variant="outline">
                      <Trash className="h-4 w-4" />
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                      <AlertDialogDescription>
                        This action cannot be undone. This will permanently delete selected
                        account and remove data from our servers.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction
                        className="bg-red-500"
                        onClick={() => deleteAccount(user.uuid)}
                      >
                        Delete
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>

    <InviteNewUserModal
      showInviteUserDialog={showInviteUserDialog}
      callback={inviteMemberCallback} />
  </>)
}