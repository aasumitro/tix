import {BrowserRouter, Route, Routes} from 'react-router-dom';
import DefaultLayout from './layouts/default';
import {HomePage} from './pages/home';
import {NotFoundPage} from './pages/not-found';
import React, {useState} from 'react';
import {EventOverviewPage} from './pages/event-overview';
import {EventParticipantPage} from './pages/event-participant';
import {UserPage} from './pages/user';

export function App() {
  const [showLoginModal, setShowLoginModal] = useState(false);

  const handleShowLoginModal = () =>  setShowLoginModal(true);

  return <>
    <BrowserRouter basename="/admin">
      <DefaultLayout showLoginModal={showLoginModal}>
        <Routes>
          <Route path="/" element={
            <HomePage unauthorizedCallback={handleShowLoginModal} />
          } />
          <Route path="/event/overview/*" element={
            <EventOverviewPage unauthorizedCallback={handleShowLoginModal} />
          } />
          <Route path="/event/participants/*" element={
            <EventParticipantPage  unauthorizedCallback={handleShowLoginModal}/>
          } />
          <Route path="/users" element={
            <UserPage unauthorizedCallback={handleShowLoginModal}/>
          } />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </DefaultLayout>
    </BrowserRouter>
  </>
}