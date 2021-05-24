import { Injectable } from '@angular/core';
import { EventBusService, PurchaseListedEvent } from '../events/event-bus.service';
import { MonitoredMomentsV09 } from '../model-storage.service';
import { Md5 } from 'ts-md5/dist/md5'

export interface ConsiderEvent {
  BlockHeight: number,
  MomentId: number,
  PlayId:  number,
  Price:  number,
  SetId:  number,
  SerialNumber: number,
  SellerAddr: string,
  isPurchase: () => boolean,
  isListing: () => boolean,

  _internal: {
    pEvent: PurchaseListedEvent
  }
}

export interface LogLevelType {
  ALERT: string,
  CONSOLE: string,
  DEBUG: string
}

export interface ConsiderAPI {
  LogLevel: LogLevelType,
  log: (level: string, evt: ConsiderEvent, msg?: string) => void,
}

export class ConsiderSandbox {

  constructor(private bus: EventBusService) {}

  compileUserCode(considerScript: string): string {
    return `
  (function sandbox(api, evt) {

    ${considerScript}

    consider(api, evt);
  })
`;
  }

  createApi(): ConsiderAPI {
    const self = this;
    return {
      LogLevel: {
        ALERT: "ALERT",
        CONSOLE: "CONSOLE",
        DEBUG: "DEBUG"
      },
      log: function (level: string, evt: ConsiderEvent, msg?: string) {
        if (level === "DEBUG") {
          console.debug("API.createApi.log", evt, msg);
        } else {
          self.bus.DetectedMoments.emit({
            Comment: msg,
            Event: evt._internal.pEvent,
            LogLevel: level
          });
        }
      }
    };
  }

  createConsiderEvent(pEvt: PurchaseListedEvent): ConsiderEvent {
    return {
      BlockHeight: pEvt.BlockHeight,
      MomentId: pEvt.MomentId,
      PlayId: pEvt.PlayId,
      Price: pEvt.Price,
      SellerAddr: pEvt.SellerAddr,
      SerialNumber: pEvt.SerialNumber,
      SetId: pEvt.SetId,
      isListing: function () {
        return pEvt.Type === "P";
      },
      isPurchase: function () {
        return pEvt.Type === "L";
      },
      _internal: {
        pEvent: pEvt
      }
    };
  }

}

@Injectable({
  providedIn: 'root'
})
export class ConsiderService {

  private sandbox: ConsiderSandbox;
  private funcCache = new Map();

  constructor(private bus: EventBusService) {
    this.sandbox = new ConsiderSandbox(bus);
  }

  onMomentEvent(models: MonitoredMomentsV09, evt: PurchaseListedEvent): void {
    models.eventMonitoring.forEach(m => {
      //if (true) {
      //  return ;
      //}

      let compiledConsider: any;
      if (this.funcCache.has(m.considerScriptMd5)) {
        compiledConsider = this.funcCache.get(m.considerScriptMd5);
      } else {
        const code = this.sandbox.compileUserCode(m.considerScript);
        try {
          compiledConsider = eval(code);
          this.funcCache.set(m.considerScriptMd5, compiledConsider);
        } catch (err) {
          console.error("Unable to compile code", code, err);
        }
      }
      try {
        compiledConsider(this.sandbox.createApi(), this.sandbox.createConsiderEvent(evt));
      } catch (err) {
        console.error("Unable to execute code", err);
      }
    });
  }


}
