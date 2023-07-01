/* eslint-disable */

import {Button} from './ui/button';
import {Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle,} from './ui/dialog';
import {format} from "date-fns"
import {Label} from './ui/label';
import {Input} from './ui/input';
import * as React from 'react';
import {useCallback, useEffect, useState} from 'react';
import {Popover, PopoverTrigger} from '@radix-ui/react-popover';
import {cn} from '../libs/utils';
import {CalendarIcon, DownloadCloud, Loader2, VerifiedIcon} from 'lucide-react';
import {PopoverContent} from './ui/popover';
import {Calendar} from './ui/calendar';
import {BaseUrl, Endpoint} from '../libs/api';
import {Badge} from './ui/badge';

interface NewEventDialogProps {
  showEventDialog: boolean,
  callback(): void
}

interface GoogleFormData {
  id: string,
  title: string,
}

export function CreateNewEventModal(props: NewEventDialogProps) {
  const [showNewEventDialog, setShowNewEventDialog] = useState(false)

  const [isErrorGoogleForm, setIErrorGoogleForm] = useState(false)
  const [googleFormIDMessage, setGoogleFormIDMessage] = useState("")
  const [isProceedGoogleForm, setIsProceedGoogleForm] = useState(false)
  const [googleFormData, setGoogleFormData] = useState<GoogleFormData[]>([])

  const [googleFormId, setGoogleFormId] = useState<string>("")
  const [name, setName] = useState<string>("")
  const [preregisterDate, setPreregisterDate] = useState<Date>()
  const [eventDate, setEventDate] = useState<Date>()
  const [location, setLocation] = useState<string>("")

  const [disabledSubmitButton, setDisabledSubmitButton] = useState(false)
  const [isSubmitted, setIsSubmitted] = useState(false)


  useEffect(() => {
    setShowNewEventDialog(props.showEventDialog)
    validateSubmitButton()
  }, [
    props.showEventDialog,
    name,
    preregisterDate,
    eventDate,
    location
  ])

  function reqSubmitNewEvent() {
    setIsSubmitted(true)
    submitNewEvent()
      .then((resp) => {
        if (resp.code === 200) {
          setGoogleFormId("")
          setGoogleFormData([])
          setIErrorGoogleForm(false)
          setShowNewEventDialog(false)
          setName("")
          setLocation("")
          props.callback()
        }

        if (resp.code === 400) {
          alert(resp.data)
        }
      })
      .finally(() => setIsSubmitted(false))
  }

  function reqValidateGoogleForm() {
    if (googleFormId === "") {
      setIErrorGoogleForm(true)
      setGoogleFormIDMessage("Google Form ID is required.")
      return
    }

    setIsProceedGoogleForm(true)
    validateGoogleForm()
      .then((resp) => {
        if (resp.code === 400) {
          setIErrorGoogleForm(true)
          setGoogleFormIDMessage(resp.data)
        }

        if (resp.code === 200) {
          setIErrorGoogleForm(false)
          setGoogleFormData(resp.data)
        }
      })
      .finally(() => setIsProceedGoogleForm(false))
  }

  const validateGoogleForm = useCallback(async () => {
    const response = await fetch(`${BaseUrl}/${Endpoint.Events.Validate}`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
      body: JSON.stringify({
        "google_form_id": googleFormId
      })
    });
    const content = await response.json();
    return Promise.resolve(content)
  }, [googleFormId])

  const validateSubmitButton = () => {
    const inputValues = [name, preregisterDate, eventDate, location];
    const isSubmitButtonDisabled = inputValues.some((value) => value === "" || value === undefined);
    setDisabledSubmitButton(isSubmitButtonDisabled);
  }

  const submitNewEvent = useCallback(async () => {
    const response = await fetch(`${BaseUrl}/${Endpoint.Events.Store}`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
      body: JSON.stringify({
        "google_form_id": googleFormId,
        "name": name,
        "preregister_date": Math.floor(preregisterDate?.getTime()! / 1000).toString(),
        "event_date": Math.floor(eventDate?.getTime()! / 1000).toString(),
        "location": location
      })
    });
    const content = await response.json();
    return Promise.resolve(content)
  }, [
    googleFormId, name, preregisterDate, eventDate, location
  ])

  return (
    <Dialog open={showNewEventDialog} onOpenChange={() => {
      setShowNewEventDialog(false)
      props.callback()
    }}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create new Event</DialogTitle>
          <DialogDescription>
            Add a new event to manage participants and tickets.
          </DialogDescription>
        </DialogHeader>
        <div>
          <div className="space-y-4 py-2 pb-4">
            <div className="space-y-2">
              <Label htmlFor="google-form-id">Google Form ID</Label>
              <div className="flex w-full items-center space-x-2">
                <Input
                  id="google-form-id"
                  placeholder="e.g: 1FAIpQLScYGWgSs8k7viWgVzLFfu1cCOpFZksIoiBFKrpMnpGSjX1jHw"
                  value={googleFormId}
                  onChange={ e => setGoogleFormId(e.currentTarget.value)}
                  disabled={googleFormData.length > 0}
                />
                <Button
                  onClick={reqValidateGoogleForm} disabled={isProceedGoogleForm || googleFormData.length > 0}
                  className={googleFormData.length > 0 ? "bg-gray-100" : ""}
                >
                  {
                    isProceedGoogleForm
                      ? <Loader2 className="h-4 w-4 animate-spin" />
                      : (googleFormData.length > 0
                        ? <VerifiedIcon  className="w-4 h-4 text-green-600" />
                        : <DownloadCloud className="w-4 h-4" />
                      )
                  }
                </Button>
              </div>
              {isErrorGoogleForm &&
                  <p className={`text-sm text-muted-foreground ${isErrorGoogleForm ? "text-red-500" : ""}`}>
                    {googleFormIDMessage}
                  </p>
              }
              {(!isErrorGoogleForm && googleFormData.length > 0) &&
                  <div className="py-2">
                    <span className="font-bold">Forms:</span>
                    {googleFormData.map((item, _) => (
                      <Badge className="ml-1" key={item.id}>{item.title.toLowerCase()}</Badge>
                    ))}
                  </div>
              }
            </div>
            {googleFormData.length > 0 && <div className="space-y-4">
              <div className="space-y-2">
                  <Label htmlFor="name">Name</Label>
                  <Input
                      id="name"
                      placeholder="e.g: Stand up commedy."
                      value={name}
                      onChange={ (e) => setName(e.currentTarget.value)}
                  />
              </div>
                <div className="space-y-2 flex flex-col">
                  <Label htmlFor="google-form-id">Preregister</Label>
                  <Popover>
                    <PopoverTrigger asChild>
                      <Button
                          variant={"outline"}
                          className={cn(
                            "w-full justify-start text-left font-normal",
                            !preregisterDate && "text-muted-foreground"
                          )}
                      >
                          <CalendarIcon className="mr-2 h-4 w-4" />
                        {preregisterDate ? format(preregisterDate, "PPP") : <span>Pick a date</span>}
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0">
                      <Calendar
                          mode="single"
                          selected={preregisterDate}
                          onSelect={(e) => setPreregisterDate(e)}
                          initialFocus
                      />
                    </PopoverContent>
                  </Popover>
                </div>
                <div className="space-y-2 flex flex-col">
                  <Label htmlFor="google-form-id">Event Date</Label>
                  <Popover>
                    <PopoverTrigger asChild>
                      <Button
                          variant={"outline"}
                          className={cn(
                            "w-full justify-start text-left font-normal",
                            !eventDate && "text-muted-foreground"
                          )}
                      >
                        <CalendarIcon className="mr-2 h-4 w-4" />
                        {eventDate ? format(eventDate, "PPP") : <span>Pick a date</span>}
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0">
                      <Calendar
                          mode="single"
                          selected={eventDate}
                          onSelect={(e) => setEventDate(e)}
                          initialFocus
                      />
                    </PopoverContent>
                  </Popover>
                </div>
                <div className="space-y-2">
                    <Label htmlFor="location">Location</Label>
                    <Input
                        id="location"
                        placeholder="e.g: Jalan Suka maju no 45"
                        value={location}
                        onChange={ (e) =>  setLocation(e.currentTarget.value)}
                    />
                </div>
            </div>}
          </div>
        </div>
        {googleFormData.length > 0 && <DialogFooter>
          <div className="flex flex-row justify-between w-full">
            <Button className="flex-0" variant="destructive" onClick={() => {
              setGoogleFormId("")
              setGoogleFormData([])
              setIErrorGoogleForm(false)
            }} disabled={googleFormData.length === 0}>Reset</Button>
            <div className="space-x-2">
              <Button
                variant="outline"
                onClick={() => {
                  setShowNewEventDialog(false)
                  props.callback()
                }}
              >Cancel</Button>
              <Button
                type="submit"
                onClick={reqSubmitNewEvent}
                disabled={disabledSubmitButton || isSubmitted}
              >
                {isSubmitted ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <></>}
                {isSubmitted ? "Please wait . . ." : "Submit"}
              </Button>
            </div>
         </div>
        </DialogFooter>}
      </DialogContent>
    </Dialog>
  )
}
