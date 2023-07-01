"use client"

import * as React from "react"
import {ChevronsUpDown, FormInputIcon, List, PlusCircle} from "lucide-react"

import {cn} from '../../libs/utils';
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar"
import { Button } from "../ui/button"
import {
  Command,
  CommandGroup,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "../ui/command"
import {
  Dialog,
  DialogTrigger,
} from "../ui/dialog"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "../ui/popover"
import {CreateNewEventModal} from '../new-event-dialog';
import {useNavigate} from 'react-router-dom';


type PopoverTriggerProps = React.ComponentPropsWithoutRef<typeof PopoverTrigger>

interface EventSwitcherProps extends PopoverTriggerProps {}

export default function EventSwitcher({ className }: EventSwitcherProps) {
  const [open, setOpen] = React.useState(false)
  const [showNewTeamDialog, setShowNewTeamDialog] = React.useState(false)

  const navigate = useNavigate()

  const eventDialogCallback  = () => {
    setShowNewTeamDialog(false)
  }

  const openGoogleForm = () => {
    const googleFormID = localStorage.getItem("current_event")
    const googleFormURL = `https://docs.google.com/forms/d/${googleFormID}/prefill`
    window.open(googleFormURL, "_blank");
  }

  return (
    <Dialog open={showNewTeamDialog} onOpenChange={setShowNewTeamDialog}>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="ghost"
            size="sm"
            role="combobox"
            aria-expanded={open}
            aria-label="Select a team"
            className={cn("w-[200px] justify-between", className)}
          >
            <Avatar className="mr-2 h-5 w-5">
              <AvatarImage
                src={`https://avatar.vercel.sh/${localStorage.getItem("current_event_name") ?? "-"}.png`}
                alt={localStorage.getItem("current_event_name") ?? "-"}
              />
              <AvatarFallback>SC</AvatarFallback>
            </Avatar>
            {localStorage.getItem("current_event_name") ?? "-"}
            <ChevronsUpDown className="ml-auto h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[200px] p-0">
          <Command>
            <CommandList>
              <CommandGroup heading="Action">
                <DialogTrigger asChild>
                  <CommandItem
                    onSelect={() => {
                      setOpen(false)
                      setShowNewTeamDialog(true)
                    }}
                  >
                    <PlusCircle className="mr-2 h-5 w-5" />
                    Create new Event
                  </CommandItem>
                </DialogTrigger>
              </CommandGroup>
            </CommandList>
            <CommandSeparator />
            <CommandList>
              <CommandGroup heading="Navigation">
                <DialogTrigger asChild>
                  <CommandItem onSelect={openGoogleForm}>
                    <FormInputIcon className="mr-2 h-5 w-5" />
                    Google Form
                  </CommandItem>
                </DialogTrigger>

                <DialogTrigger asChild>
                  <CommandItem onSelect={() => navigate("/")}>
                    <List className="mr-2 h-5 w-5" />
                    Show Event Lists
                  </CommandItem>
                </DialogTrigger>
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      <CreateNewEventModal showEventDialog={showNewTeamDialog} callback={eventDialogCallback} />
    </Dialog>
  )
}