import {Button} from './ui/button';
import errorImage from '../assets/img/error.png';

interface ErrorProps {
  code?: number;
  message?: string;
  callback: () => void;
}

// #ffffff, #8f8f8f, #bdc2c1, #cfcfcf, #ffffff, #ffffff = designstripe

export function ErrorSection(props: ErrorProps) {
  let r = (Math.random() + 1).toString(36).substring(7);
  return (<div className="flex flex-col items-center py-24">
    <img src={errorImage} alt={"error-image"+r} width={350}/>
    <h1 className="text-2xl font-bold">Oops . . .</h1>
    <p className="text-lg font-extralight mt-2">Something went wrong</p>
    <Button variant="ghost" className="mt-4" onClick={props.callback}>Refresh</Button>
  </div>)
}