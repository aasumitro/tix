/* eslint-disable */

import * as React from 'react';
import {useEffect, useState} from 'react';
import {ErrorSection} from '../components/error-section';
import {UserListSkeleton} from '../components/user-list-skeleton';
import {BaseUrl, Endpoint} from '../libs/api';
import {UserListData} from '../components/user-list-data';

interface UserPageProps {
  unauthorizedCallback: () => void
}

export function UserPage(props: UserPageProps) {

  const [isLoading, setIsLoading] = useState(false)
  const [isError, setIsError] = useState(false)
  const [users, setUsers] = useState(null)

  useEffect( () => fetchUsers(), [props])

  function fetchUsers() {
    setIsLoading(true)
    setIsError(false)
    fetch(`${BaseUrl}/${Endpoint.User.List}`, {
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
      .then(resp => setUsers(resp.data))
      .catch(_ => setIsError(true))
      .finally(()=> setIsLoading(false))
  }

  const errorCallback = () => fetchUsers()

  return <>
    {isError && <ErrorSection callback={errorCallback} />}

    {isLoading && <UserListSkeleton />}

    {(!isLoading && !isError && users !== null) &&
        <UserListData
            users={users}
            doRefreshCallback={fetchUsers}
        />
    }
  </>
}