/* eslint-disable */

import * as React from 'react';
import {EventParticipantData} from '../components/event-participant-data';
import {EventParticipantSkeleton} from '../components/event-participant-skeleton';
import {ErrorSection} from '../components/error-section';
import {useEffect, useState} from 'react';
import {BaseUrl, Endpoint} from '../libs/api';
import {NoDataSection} from '../components/no-data-section';
import {toast} from '../components/ui/use-toast';
import {useNavigate} from 'react-router-dom';

interface EventParticipantPageProps {
  unauthorizedCallback: () => void
}

export function EventParticipantPage(props: EventParticipantPageProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [isError, setIsError] = useState(false)
  const [participants, setParticipants] = useState(null)

  const [isSyncProceed, setIsSyncProceed] = React.useState(false)
  const [buttonSyncText, setButtonSyncText] = React.useState("Sync participant data")


  useEffect( () => {
    fetchEventParticipants()
  }, [props])

  function fetchEventParticipants() {
    setIsLoading(true)
    setIsError(false)
    const eventID = localStorage.getItem("current_event")
    fetch(`${BaseUrl}/${Endpoint.Events.Participants(eventID as string)}`, {
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
      .then(resp => setParticipants(resp.data))
      .catch(error => setIsError(true))
      .finally(()=> setIsLoading(false))
  }

  const navigate = useNavigate()

  function handlerSyncParticipantData() {
    setIsSyncProceed(true)
    setButtonSyncText("Please wait . . .")
    const eventID = localStorage.getItem("current_event")
    if (eventID === "") {
      toast({
        variant: "destructive",
        title: "Failed",
        description: `Please select an event to do this action.`,
      })
      navigate('/')
      return
    }
    fetch(`${BaseUrl}/${Endpoint.Events.Sync(eventID as string)}`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
    }).then((resp) => {
      if (resp.status === 401) {
        toast({
          variant: "destructive",
          title: "Unauthenticated",
          description: "Please login to continue.",
        })
        navigate('/')
        return
      }

      return resp.json();
    }).then((resp) => {
      setIsSyncProceed(false)
      setButtonSyncText("Sync participant data")

      if (!resp)  {
        return;
      }

      toast({
        variant: resp.code === 200 ? "default": "destructive",
        title: "Action Sync Data",
        description: resp.data,
      })
    })
  }

  const errorCallback = () => fetchEventParticipants()

  return <>
    <div className="container w-full h-full space-y-4 py-12 flex flex-col">
      {isError && <ErrorSection callback={errorCallback} />}

      {isLoading &&  <EventParticipantSkeleton /> }

      {(!isLoading && !isError && participants === null) &&
          <NoDataSection
              buttonTitle={buttonSyncText}
              dataName={"Participants"}
              callback={handlerSyncParticipantData}
          />
      }

      {(!isLoading && !isError && participants !== null) &&
          <EventParticipantData
              participants={participants}
              doRefreshCallback={fetchEventParticipants}
              isSyncProceed={isSyncProceed}
              buttonSyncText={buttonSyncText}
              buttonSyncCallback={handlerSyncParticipantData}
          />
      }
    </div>
  </>
}