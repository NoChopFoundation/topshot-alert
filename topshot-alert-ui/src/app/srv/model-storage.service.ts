import { Injectable } from '@angular/core';

export interface MonitoredItem {
  playId: number;
  setId: number;

  considerScript: string;
  considerScriptMd5: string
}

export interface MonitoredMomentsV09 {
  eventMonitoring: MonitoredItem[]
}


@Injectable({
  providedIn: 'root'
})
export class ModelStorageService {

  constructor() { }

  loadMonitoredMoments(): MonitoredMomentsV09 {
    const model = window.localStorage.getItem("MonitoredMomentsV09");
    if (model === null) {
      return <MonitoredMomentsV09> {
        eventMonitoring: []
      };
    } else {
      return JSON.parse(model);
    }
  }

  importModel(importedModelStr: string): boolean {
    try {
      const importedJson = JSON.parse(importedModelStr);
      if (importedJson.MonitoredMomentsV09) {
        const importedModel = <MonitoredMomentsV09>importedJson.MonitoredMomentsV09;
        const currentSize = this.loadMonitoredMoments().eventMonitoring.length;
        if (currentSize > 0) {
          const answer = confirm(`About to overwrite ${currentSize} monitored with ${importedModel.eventMonitoring.length} imported items.  Import and overwrite existing?`);
          if (answer) {
            this.persistMonitoredMoments(importedModel);
            return true;
          }
        }
      } else {
        alert('Imported model not recongized');
      }
    } catch(e) {
      alert('Unable to JSON parse model, ensure it is JSON format ' + e);
    }
    return false;
  }

  exportModel(): string {
    return JSON.stringify({
      MonitoredMomentsV09: this.loadMonitoredMoments(),
      exportedMetaData: {
        time: new Date().toISOString()
      }
    }, null, 2);
  }

  persistMonitoredMoments(model: MonitoredMomentsV09) {
    window.localStorage.setItem("MonitoredMomentsV09", JSON.stringify(model));
  }

}
