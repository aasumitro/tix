import {Event} from "../libs/model/event"
import {Calendar, PinIcon, Plus, User} from 'lucide-react';
import {Card, CardContent, CardFooter, CardHeader, CardTitle} from './ui/card';
import * as React from 'react';
import {useNavigate} from 'react-router-dom';
import {Badge} from './ui/badge';
import {HoverCard, HoverCardContent, HoverCardTrigger} from './ui/hover-card';

interface EventListSectionProps {
  events: Event[];
  action: () => void;
}

export function EventListData(props: EventListSectionProps) {
  const navigate = useNavigate()

  const goToEvent = (id: string, name: string) => {
    localStorage.setItem("current_event", id)
    localStorage.setItem("current_event_name", name)
    navigate(`/event/overview/${id}`)
  }

  const eventDate = (date: number): string => {
    const milliseconds = date * 1000;
    const specifiedDate = new Date(milliseconds);
    const day = specifiedDate.getDay()
    const month = specifiedDate.toLocaleString('default', { month: 'long' });
    const year = specifiedDate.getFullYear();
    return `${day} ${month} ${year}`
  }

  const dayDiff = (date: number): number => {
    const milliseconds = date * 1000;
    const currentDate = new Date();
    const specifiedDate = new Date(milliseconds);
    const differenceMs = Math.abs(currentDate.getTime() - specifiedDate.getTime());
    return Math.floor(differenceMs / (1000 * 60 * 60 * 24));
  }

  return (
    <div className="flex flex-wrap gap-4 py-24">
      <button className="w-60 h-60 bg-gray-50 hover:bg-gray-100 border-2 border-dashed border-gray-200 rounded-lg" onClick={props.action}>
        <Plus className="w-12 h-12 text-gray-400 mx-auto"/>
      </button>
      {props?.events?.map((event, index) => (
        <HoverCard key={event.id}>
          <HoverCardTrigger asChild>
            <button
              className="w-60 h-60 rounded-lg"
              onClick={() => goToEvent(event.google_form_id, event.name)}
              disabled={dayDiff(event.event_date) <= 0}
            >
              <Card className="w-full h-full text-left hover:bg-gray-100 flex flex-col justify-between">
                {event.is_active && <Badge className="absolute text-right p-1 ml-0 bg-red-500 animate-ping"></Badge>}
                <CardHeader className="space-y-0 pb-2">
                  <CardTitle className="text-lg font-medium">
                    #{index+1} - {event.name}
                  </CardTitle>
                </CardHeader>
                <CardContent className="mt-4">
                  <div className="flex gap-2 items-center">
                    <User className="h-6 w-6"/>
                    <div className="text-2xl font-bold">{event.total_participants ?? 0}</div>
                  </div>
                  <p className="text-md ml-1 mt-2 text-muted-foreground">
                    Participants
                  </p>
                </CardContent>
                <CardFooter className="text-sm font-normal flex items-baseline flex-row gap-2">
                  <div className="flex gap-2 items-center text-md">
                    {dayDiff(event.event_date)} days before the event.
                  </div>
                </CardFooter>
              </Card>
            </button>
          </HoverCardTrigger>
          <HoverCardContent className="w-80 mt-4">
            <div className="flex flex-col text-left space-y-4">
              <div className="flex space-x-2 items-center">
                <div className="p-2">
                  <Calendar className="h-6 w-6" />
                </div>
                <p>{eventDate(event.event_date)}</p>
              </div>
             <div className="flex space-x-2 items-center">
               <div className="p-2">
                <PinIcon className="h-6 w-6" />
               </div>
               <p>{event.location}</p>
             </div>
            </div>
          </HoverCardContent>
        </HoverCard>
      ))}
    </div>
  )
}