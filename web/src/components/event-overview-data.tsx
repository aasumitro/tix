import {Card, CardContent, CardDescription, CardHeader, CardTitle} from './ui/card';
import {Users} from 'lucide-react';
import {Bar, BarChart, ResponsiveContainer, XAxis, YAxis} from 'recharts';
import {Avatar, AvatarFallback} from './ui/avatar';
import * as React from 'react';
import {Event} from "../libs/model/event"

interface EventOverviewDataProps {
  event: Event
}

export function EventOverviewData(props: EventOverviewDataProps) {
  return (<>
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            Total Responded
          </CardTitle>
          <Users className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {props?.event?.total_participants ?? 0}
          </div>
          <p className="text-xs text-muted-foreground mt-2">
            All participants of the event make a response through Google Forms, indicating their interest and intent to join the event.
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            Approved
          </CardTitle>
          <Users className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {props?.event?.total_approved_participant ?? 0}
          </div>
          <p className="text-xs text-muted-foreground mt-2">
            Total participants who have successfully received their generated tickets, which grant them permission to attend the event.
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Waiting Approval</CardTitle>
          <Users className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {props?.event?.total_waiting_approval_participant ?? 0}
          </div>
          <p className="text-xs text-muted-foreground mt-2">
            There are certain participants who are currently awaiting approval in order to receive their generated tickets.
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            Declined
          </CardTitle>
          <Users className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {props?.event?.total_declined_participant ?? 0}
          </div>
          <p className="text-xs text-muted-foreground mt-2">
            Some participants were declined for fraud or being underage based on PoP (Proof of Payment) and DoB (Date of Birth).
          </p>
        </CardContent>
      </Card>
    </div>

    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <Card className="col-span-4">
        <CardHeader>
          <CardTitle>Respondents overview of the past 7 days.</CardTitle>
        </CardHeader>
        <CardContent className="pl-2">
          <ResponsiveContainer width="100%" height={350}>
            <BarChart data={props?.event?.weekly_overview}>
              <XAxis
                dataKey="name"
                stroke="#888888"
                fontSize={12}
                tickLine={false}
                axisLine={false}
              />
              <YAxis
                stroke="#888888"
                fontSize={12}
                tickLine={false}
                axisLine={false}
                tickFormatter={(value: any) => `${value}`}
              />
              <Bar dataKey="total" fill="#adfa1d" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>
      <Card className="col-span-3">
        <CardHeader>
          <CardTitle>Recent respondents</CardTitle>
          <CardDescription>
            You got {props?.event?.latest_respondents?.length ?? 0} respondents today.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {props?.event?.latest_respondents?.length
              ?  props?.event?.latest_respondents?.map((respondent, index) => (
              <button key={index} className="flex text-left w-full items-center hover:bg-gray-50 px-2 py-4 rounded-md">
                <Avatar className="h-9 w-9">
                  <AvatarFallback>
                    {respondent.name.substring(0, 2)}
                  </AvatarFallback>
                </Avatar>
                <div className="ml-4 space-y-1">
                  <p className="text-sm font-medium leading-none">
                    {respondent.name} <span className="text-muted-foreground"></span>
                  </p>
                  <p className="text-sm text-muted-foreground">
                    <span>{respondent.email}</span> | <span>{respondent.phone}</span>
                  </p>
                </div>
              </button>
            ))
              : <div className="text-center text-muted-foreground mt-8">No recent respondents.</div>
            }
          </div>
        </CardContent>
      </Card>
    </div>
  </>)
}