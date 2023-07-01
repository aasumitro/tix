import {Button} from './ui/button';
import noDataImage from '../assets/img/no-data.png';

interface NoDataProps {
  buttonTitle: string;
  dataName: string;
  callback: () => void;
}

// #ffffff, #8f8f8f, #bdc2c1, #cfcfcf, #ffffff, #ffffff = designstripe

export function NoDataSection(props: NoDataProps) {
  return (<div className="flex flex-col items-center py-24">
    <img src={noDataImage} alt={`no-data-${props.dataName}`} width={350}/>
    <h1 className="text-2xl font-bold">Uowhh . . .</h1>
    <p className="text-lg font-extralight mt-2">Seems like you don't have any {props.dataName} yet</p>
    <Button variant="ghost" className="mt-4" onClick={props.callback}>{props.buttonTitle}</Button>
  </div>)
}