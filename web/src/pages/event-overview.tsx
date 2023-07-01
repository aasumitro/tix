/* eslint-disable */

import {useEffect, useState} from 'react';
import {ErrorSection} from '../components/error-section';
import * as React from 'react';
import {EventOverviewSkeleton} from '../components/event-overview-skeleton';
import {EventOverviewData} from '../components/event-overview-data';
import {BaseUrl, Endpoint} from '../libs/api';

interface EventOverviewPageProps {
  unauthorizedCallback: () => void
}

export function EventOverviewPage(props: EventOverviewPageProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [isError, setIsError] = useState(false)
  const [events, setEvents] = useState(null)

  useEffect( () => {
    fetchEventOverview()
  }, [props])

  function fetchEventOverview() {
    setIsLoading(true)
    setIsError(false)
    const eventID = localStorage.getItem("current_event")
    fetch(`${BaseUrl}/${Endpoint.Events.Overview(eventID as string)}`, {
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
      .then(resp => setEvents(resp.data))
      .catch(error => setIsError(true))
      .finally(()=> setIsLoading(false))
  }

  const errorCallback = () => fetchEventOverview()

  return (<>
    <div className="container w-full h-full space-y-4 p-8">
      {isError && <ErrorSection callback={errorCallback} />}

      {isLoading && <EventOverviewSkeleton />}

      {(!isLoading && !isError && events !== null) && <EventOverviewData event={events} />}
    </div>
  </>)
}