"use client"

import { ColumnDef } from "@tanstack/react-table"
import {ArrowUpDown, ExternalLink, MoreHorizontal} from "lucide-react"

import { Button } from "../ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel, DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import * as React from 'react';
import {Participant} from '../../libs/model/participant';
import {Badge} from '../ui/badge';
import {toast} from '../ui/use-toast';
import {ActionStatus} from './action-status';
import {Dialog, DialogTrigger} from "../ui/dialog"
import {BaseUrl, Endpoint} from '../../libs/api';


export const columns: ColumnDef<Participant>[] = [
  {
    accessorKey: "name",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Name
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
    cell: ({ row }) => <div className="capitalize">{row.getValue("name")}</div>,
  },
  {
    accessorKey: "email",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Email
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
    cell: ({ row }) => <div className="lowercase">{row.getValue("email")}</div>,
  },
  {
    accessorKey: "phone",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Phone
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
    cell: ({ row }) => <div className="lowercase">{row.getValue("phone")}</div>,
  },
  {
    accessorKey: "job",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Job
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
    cell: ({ row }) => <div className="lowercase">{row.getValue("job")}</div>,
  },
  {
    accessorKey: "date_of_birth",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Date of Birth
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
    cell: ({ row }) => <div className="lowercase">{row.getValue("date_of_birth")}</div>,
  },
  {
    accessorKey: "prof_of_payment",
    header: "Proof of Payment",
    cell: ({ row }) => <div className="lowercase">
      <a
        className="text-blue-400  flex gap-2"
        href={row.getValue("prof_of_payment")}
        target="_blank" rel="noreferrer"
      >
        <ExternalLink className="w-4 h-4"/>
        Click to see
      </a>
    </div>,
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => (
      <Badge variant="outline" className="capitalize">
        {row.getValue("status")}
      </Badge>
    ),
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const data = row.original
      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {row.getValue("status") === "approved" && <DropdownMenuItem
                onClick={() => resendTicket(data)}
            >
                Resend Ticket
            </DropdownMenuItem>}
            {row.getValue("status") === "waiting approval" &&  <Dialog>
              <DialogTrigger asChild>
                  <button className="p-2 text-sm border-0 w-full hover:bg-gray-100 rounded-md">
                      Update Status
                  </button>
              </DialogTrigger>
              <ActionStatus participant={data}/>
            </Dialog>}
            {row.getValue("status") === "declined" && <DropdownMenuItem
                onClick={() => viewDetailTicket(data)}
                disabled
            >
                View detail
            </DropdownMenuItem>}
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
  },
]


const resendTicket = (participant: Participant) => {
  const eventID = localStorage.getItem("current_event")
  generateTicket(eventID as string, participant.id).then((resp) => {
    if (resp.code === 401) {
      localStorage.removeItem("current_event")
      localStorage.removeItem("is_login")
      window.location.href = "/admin"
      return
    }

    if (resp.code === 400) {
      toast({
        variant: "destructive",
        title: "Action",
        description: resp.data,
      })
    }

    if (resp.code === 200) {
      toast({
        variant: "default",
        title: "Action",
        description: `${resp.data} we will sent the ticket to ${participant.email} also your email.`,
      })
    }
  })
};

const generateTicket = async (eventID: string, participantID: number) => {
  const response = await fetch(`${BaseUrl}/${Endpoint.Events.GenerateTicket(eventID, participantID)}`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    credentials: 'include',
    body: JSON.stringify({})
  });
  const content = await response.json();
  return Promise.resolve(content)
};

const viewDetailTicket = (participant: Participant) => console.log(participant);