import * as React from 'react'
import {Toaster} from '../components/toaster';
import {Header} from '../components/header';
import {Footer} from '../components/footer';
import {SplashScreen} from '../components/splash';
import {useEffect, useState} from 'react';
import {LoginModal} from '../components/login-modal';
import {useLocation, useNavigate} from 'react-router-dom';
import {BaseUrl, Endpoint} from '../libs/api';
import {User} from '../libs/model/user';

interface Props {
  showLoginModal?: boolean
  children: React.ReactNode
}

const DefaultLayout: React.FC<Props> = (props: Props) => {
  const [showLoginDialog, setShowLoginDialog] = React.useState(false)
  const [showNavs, setShowNavs] = useState(false);
  const [showSplashScreen, setShowSplashScreen] = useState(true);
  const [user, setUser] = useState<User>();
  const location = useLocation()
  const navigate = useNavigate()

  useEffect(() => {
    const { cookie } = window.document;
    if (cookie === "") return
    const url = window.location.href;
    const searchParams = new URLSearchParams(url.split('#')[1]);
    const accessToken = searchParams.get('access_token');
    const accessType = searchParams.get('type');
    if (accessToken) {
      saveSession(accessToken as string, accessType as string).then((resp) => {
        setShowLoginDialog(resp.code !== 201)
        if (resp.code === 201) {
          localStorage.setItem("is_login", "true")
        }
      })
    }
    const currentEvent = localStorage.getItem("current_event")
    if (currentEvent === null) navigate('/')
    setShowNavs(location.pathname !== '/')
    const isLoggedIn = localStorage.getItem("is_login")
    if (!isLoggedIn || props.showLoginModal) setShowLoginDialog(true)
    setTimeout(() => {
      setShowSplashScreen(false);
    }, 1000);
    if (isLoggedIn) getProfile()
  }, [location, navigate, props.showLoginModal])

  async function saveSession(jwt: string, type: string) {
    const response = await fetch(`${BaseUrl}/${Endpoint.Auth.Verify}`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
      body: JSON.stringify({jwt, type})
    });
    const content = await response.json();
    return Promise.resolve(content)
  }

  function getProfile() {
    fetch(`${BaseUrl}/${Endpoint.Auth.Profile}`, {
      method: 'GET',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include'
    }).then((resp) => {
      if (resp.status === 401) {
        setShowLoginDialog(true)
        return
      }
      return resp.json()
    }).then((resp) => setUser(resp.data))
  }

  if (showSplashScreen) {
    return <SplashScreen />;
  }

  return (
    <>
      <div className="md:hidden text-center py-10">
        Screen Size Not Supported <br/> (min: 768px/tablet screen)
      </div>

      <div className="hidden flex-col md:flex mb-16">
        {(showNavs) && <Header user={user} />}
          <main className="min-h-full min-w-full">
            {props.children}
          </main>
          <Toaster />
        {(showNavs) && <Footer />}
      </div>

      <LoginModal showLoginDialog={showLoginDialog}/>
    </>
  )
}

export default DefaultLayout;