import {Button} from './ui/button';
import {RefreshCcwIcon} from 'lucide-react';
import * as React from 'react';
import {Participant} from '../libs/model/participant';
import {columns} from './participant-data-table/columns';
import {DataTable} from './participant-data-table/data-table';
import {useEffect} from 'react';
import {toast} from './ui/use-toast';
import {BaseUrl, Endpoint} from '../libs/api';
import {useNavigate} from 'react-router-dom';
import {EventExportType} from '../libs/enums/event-export-type';

interface EventParticipantDataProps {
  participants: Participant[];
  doRefreshCallback: () => void;
  isSyncProceed: boolean,
  buttonSyncText: string,
  buttonSyncCallback: () => void,
}

export function EventParticipantData(props: EventParticipantDataProps) {
  const [isProceed, setIsProceed] = React.useState(false)
  const [buttonText, setButtonText] = React.useState("")

  const navigate = useNavigate()

  useEffect(() => {
    setButtonText("Sync participant data")
    setIsProceed(props.isSyncProceed)
    setButtonText(props.buttonSyncText)
  }, [props])


  function handleExportEventData(exportType: EventExportType) {
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
    fetch(`${BaseUrl}/${Endpoint.Events.Export(eventID as string, exportType)}`, {
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
      if (!resp)  {
        return;
      }

      toast({
        variant: resp.code === 200 ? "default": "destructive",
        title: "Action Export Data",
        description: resp.data,
      })
    })
  }

  return (
    <>
      <div className="flex flex-row items-center justify-between mb-10">
        <div>
          <h1 className="text-2xl font-bold">Participants</h1>
          <p>List of event participants who request to get their tickets.</p>
        </div>
        <Button onClick={props.buttonSyncCallback} disabled={isProceed}>
          <RefreshCcwIcon className={`mr-2 h-4 w-4 ${isProceed ? "animate-spin" : ""}`} /> {buttonText}
        </Button>
      </div>

      <DataTable
        columns={columns}
        data={props.participants}
        exportData={handleExportEventData}
      />
    </>
  )
}