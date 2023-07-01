import {
  AlertDialog, AlertDialogAction, AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from './ui/alert-dialog';
import {Label} from './ui/label';
import {Input} from './ui/input';
import * as React from 'react';
import {useCallback, useEffect} from 'react';
import {Loader2} from 'lucide-react';
import {BaseUrl, Endpoint} from '../libs/api';
import {toast} from './ui/use-toast';

interface InviteModalProps {
  showInviteUserDialog: boolean,
  callback(): void,
}
export function InviteNewUserModal(props: InviteModalProps) {
  const [showInviteUserDialog, setShowInviteUserDialog] = React.useState(false)
  const [isProceed, setIsProceed] = React.useState(false)
  const [email, setEmail] = React.useState("")
  const [isError, setIsError] = React.useState(false)
  const [formMessage, setFormMessage] = React.useState("")
  const [buttonText, setButtonText] = React.useState("")

  useEffect(() => {
    setIsError(false)
    setIsProceed(false)
    setFormMessage("Enter email address.")
    setButtonText("Invite User")
    setShowInviteUserDialog(props.showInviteUserDialog)
  }, [props.showInviteUserDialog])

  function doInviteMember() {
    if (email === "") {
      setIsError(true)
      setFormMessage("Email is required.")
      return
    }

    setIsProceed(true)
    setButtonText("please wait...")
    inviteMember().then((resp) => {
      if (resp.code === 200) {
        toast({
          variant: "default",
          title: "Success",
          description: `We've sent an invitation to ${email}.`,
        })
        setShowInviteUserDialog(false)
        props.callback()
      }

      if (resp.code === 422) {
        setIsError(true)
        setFormMessage(resp.data)
      }

      if (resp.code === 500) {
        setIsError(true)
        setFormMessage(resp.data)
      }
    }).finally(() => {
        setIsProceed(false)
        setButtonText("Invite User")
    })
  }

  const inviteMember = useCallback(async() => {
    const response = await fetch(`${BaseUrl}/${Endpoint.User.Invite}`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
      body: JSON.stringify({email})
    });
    const content = await response.json();
    return Promise.resolve(content)
  }, [email]);

  return <>
    <AlertDialog open={showInviteUserDialog}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Invite new User</AlertDialogTitle>
          <AlertDialogDescription>
            You can add admin by invite them via email.
          </AlertDialogDescription>
          <div>
            <div className="space-y-4 py-2 pb-4">
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  placeholder="admin@tix.id"
                  value={email}
                  onChange={ e => setEmail(e.currentTarget.value)}/>
                <p className={`text-sm text-muted-foreground ${isError ? "text-red-500" : ""}`}>{formMessage}</p>
              </div>
            </div>
          </div>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={() => {
            setShowInviteUserDialog(false)
            setEmail("")
            props.callback()
          }}>
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction onClick={doInviteMember} disabled={isProceed}>
            {isProceed ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <></>}
            {buttonText}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </>
}
