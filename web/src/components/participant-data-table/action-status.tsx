import {Participant} from '../../libs/model/participant';
import {Select} from '@radix-ui/react-select';
import {SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue} from '../ui/select';
import {Button} from '../ui/button';
import * as React from 'react';
import {useState} from 'react';
import {toast} from '../ui/use-toast';
import {Label} from '../ui/label';
import {
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '../ui/dialog';
import {Textarea} from '../ui/textarea';
import {BaseUrl, Endpoint} from '../../libs/api';
import {Loader2} from 'lucide-react';

interface ActionStatusProps {
  participant: Participant;
}

export function ActionStatus(props: ActionStatusProps) {
  const [status, setStatus] = useState("")
  const [reason, setReason] = useState("")
  const [isProceed, setIsProceed] = React.useState(false)
  const [buttonText, setButtonText] = React.useState("Save changes")

  const updateStatusTicket = (participant: Participant) => {
    if (status === "") {
      toast({
        variant: "destructive",
        title: "Action",
        description: `Please pick an status to do this action.`,
      })
      return
    }

    if (status === "declined" && reason === "") {
      toast({
        variant: "destructive",
        title: "Action",
        description: `Please provide declined reason.`,
      })
      return
    }

    setIsProceed(true)
    setButtonText("Please wait...")
    const eventID = localStorage.getItem("current_event")
    updateRespondedStatus(eventID as string, participant.id).then((resp) => {
      setIsProceed(false)
      setButtonText("Save changes")

      if (resp.code === 401) {
        localStorage.removeItem("current_event")
        localStorage.removeItem("is_login")
        window.location.href = "/admin"
        return
      }

      if (resp.code === 400 || resp.code === 422) {
        toast({
          variant: "destructive",
          title: "Action",
          description: resp.data,
        })
      }

      if (resp.code === 200) {
        setStatus("")
        setReason("")
        toast({
          variant: "default",
          title: "Action",
          description: resp.data,
        })
        const event = new KeyboardEvent('keydown', {
          key: 'Escape',
          code: 'Escape',
          bubbles: true
        });
        document.dispatchEvent(event);
      }
    })
  };

  const updateRespondedStatus = async (eventID: string, participantID: number) => {
    const response = await fetch(`${BaseUrl}/${Endpoint.Events.Status(eventID, participantID)}`, {
      method: 'PATCH',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
      body: JSON.stringify({
        status: status,
        declined_reason: reason
      })
    });
    const content = await response.json();
    return Promise.resolve(content)
  };

  return (
    <>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Update Status</DialogTitle>
          <DialogDescription>
            Make changes to your profile here. Click save when you're done.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <Select onValueChange={setStatus}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder="Select a status" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Status</SelectLabel>
                <SelectItem value="approved">Approved</SelectItem>
                <SelectItem value="declined">Declined</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
          {status === "declined" && <div className="grid w-full max-w-sm items-center gap-2">
              <Label htmlFor="reason">Reason</Label>
              <Textarea
                  id="reason"
                  placeholder="e.g: fraud detected . . ."
                  value={reason}
                  onChange={ e => setReason(e.currentTarget.value)}
              />
          </div>}
        </div>
        <DialogFooter>
          <Button type="submit" onClick={() => updateStatusTicket(props.participant)} disabled={isProceed}>
            {isProceed ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <></>}
            {buttonText}
          </Button>
        </DialogFooter>
      </DialogContent>
    </>
  )
}