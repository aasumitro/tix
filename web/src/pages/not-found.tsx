import {useNavigate} from 'react-router-dom';

export function NotFoundPage() {
  const navigate = useNavigate()

  const goHome = () => {
    localStorage.removeItem("current_event")
    navigate("/")
  }

  return (<>
    <div className="mx-auto text-center mt-64">
      <h1 className="text-8xl font-semibold">404</h1>
      <p className="mt-4 text-2xl font-extralight">Page not found!</p>
      <button onClick={goHome} className="mt-4">Back Home</button>
    </div>
  </>)
}