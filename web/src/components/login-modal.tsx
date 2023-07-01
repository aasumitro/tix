import {
  AlertDialog, AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from './ui/alert-dialog';
import {Label} from './ui/label';
import {Input} from './ui/input';
import * as React from 'react';
import {useCallback, useEffect} from 'react';
import {BaseUrl, Endpoint, ErrorMessage} from '../libs/api';
import {Loader2} from 'lucide-react';

interface LoginModalProps {
  showLoginDialog: boolean,
}

export function LoginModal(props: LoginModalProps) {
  const [showLoginDialog, setShowLoginDialog] = React.useState(false)
  const [showInput, setShowInput] = React.useState(false)
  const [isProceed, setIsProceed] = React.useState(false)
  const [email, setEmail] = React.useState("")
  const [isError, setIsError] = React.useState(false)
  const [formMessage, setFormMessage] = React.useState("")
  const [buttonText, setButtonText] = React.useState("")

  useEffect(() => {
    setIsError(false)
    setIsProceed(false)
    setShowInput(true)
    setShowLoginDialog(props.showLoginDialog)
    setFormMessage("Enter your email address.")
    setButtonText("Request Magic Link")
  }, [props.showLoginDialog])

  function reqMagicLink() {
    if (email === "") {
      setIsError(true)
      setFormMessage("Email is required.")
      return
    }

    setIsProceed(true)
    setButtonText("Please wait...")
    checkSession().then((resp) => {
      if (resp.code === 200) {
        setShowInput(false)
      }

      if (resp.code === 429) {
        setIsError(true)
        setFormMessage(resp.data)
        if (resp.data === ErrorMessage.Auth.ToManyRequest) {
          let waitingTime = 60
          const interval   = setInterval(() => {
            setButtonText(`Please wait (${waitingTime})`)
            waitingTime -= 1
          }, 1000)
          setTimeout(() => {
            setIsProceed(false)
            setButtonText("Resend Magic Link")
            clearInterval(interval)
          }, 1000 * 60)
        }
      }

       if (resp.code === 500) {
         setIsProceed(false)
         setButtonText("Resend Magic Link")
        setIsError(true)
        setFormMessage(resp.data)
       }
    })
  }

  const checkSession = useCallback(async() => {
    const response = await fetch(`${BaseUrl}/${Endpoint.Auth.Validate}`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
      body: JSON.stringify({email})
    });
    const content = await response.json();
    return Promise.resolve(content)
  }, [email]);


  return <>
    <AlertDialog open={showLoginDialog}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Login</AlertDialogTitle>
          <AlertDialogDescription>
            Sign in to your account to continue
          </AlertDialogDescription>
          {showInput && <div>
            <div className="space-y-4 py-2 pb-4 text-left">
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
          </div>}

          {!showInput && <div className="py-2 pb-4 text-left">
              We've sent you a magic link to your email. Please check your inbox.
          </div>}
        </AlertDialogHeader>
        <AlertDialogFooter>
          {showInput && <AlertDialogAction onClick={reqMagicLink} disabled={isProceed}>
            {isProceed ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <></>}
            {buttonText}
          </AlertDialogAction>}
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </>
}
