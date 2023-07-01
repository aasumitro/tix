/* eslint-disable */

import {CreateNewEventModal} from '../components/new-event-dialog';
import * as React from 'react';
import { useEffect, useState} from 'react';
import {BaseUrl, Endpoint} from '../libs/api';
import {EventListSkeleton} from '../components/event-list-skeleton';
import {ErrorSection} from '../components/error-section';
import {NoDataSection} from '../components/no-data-section';
import {EventListData} from '../components/event-list-data';

interface HomePageProps {
  unauthorizedCallback: () => void
}

export function HomePage(props: HomePageProps) {
  const [showNewEventDialog, setShowNewEventDialog] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [isError, setIsError] = useState(false)
  const [events, setEvent] = useState(null)

  const createNewEvent = () => setShowNewEventDialog(true)

  const eventDialogCallback  = () => {
    setShowNewEventDialog(false)
    fetchEvents()
  }

  useEffect( () => {
    fetchEvents()
  }, [])

  function fetchEvents() {
    setIsLoading(true)
    setIsError(false)
    fetch(`${BaseUrl}/${Endpoint.Events.List}`, {
      method: 'GET',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
    })
      .then((resp) => {
        if (resp.status === 401) {
          props.unauthorizedCallback()
          setIsError(true)
          return
        }
        return resp.json()
      })
      .then(resp => setEvent(resp.data))
      .catch(_ => setIsError(true))
      .finally(()=> setIsLoading(false))
  }

  const errorCallback = () => fetchEvents()

  return (<>
    <div className="container w-full h-full space-y-4 p-8">
      <div className="text-center">
        <h1 className="text-4xl font-semibold"> - TIX - </h1>
        <p className="mt-4 text-lg font-extralight">Manage Events Participants and Tickets</p>
      </div>

      {isError && <ErrorSection callback={errorCallback} />}

      {isLoading && <EventListSkeleton display={10}/>}

      {(!isLoading && !isError && events === null) &&
          <NoDataSection
              buttonTitle={"Create new Event"}
              dataName={"events"}
              callback={() => setShowNewEventDialog(true)}
          />
      }

      {(!isLoading && !isError && events !== null) &&
          <EventListData
              events={events}
              action={createNewEvent}
          />
      }

      <CreateNewEventModal
        showEventDialog={showNewEventDialog}
        callback={eventDialogCallback}
      />
    </div>
  </>);
}