import { HttpClient } from '@angular/common/http';
import { EventEmitter, Injectable } from '@angular/core';
import { BehaviorSubject, Observable, Subscriber, Subscription, TeardownLogic } from 'rxjs';
import { AppConfigService } from '../app-config.service';

export interface PurchaseListedEventPayload {
  Data: PurchaseListedEvent[]
}

export interface DetectedConsiderEvent {
  LogLevel: string
  Comment: string | undefined
  Event: PurchaseListedEvent
}

export interface PurchaseListedEvent {
  Type: string;
  MomentId: number;
  BlockHeight: number;
  PlayId: number;
  SerialNumber: number;
  SetId: number;
  SellerAddr: string;
  Price: number;
  Created_At: string; // "2021-04-02T16:15:39Z"
}

export class ContactedEvent {
  private contactTime: number
  constructor(initTimeMs?: number) {
    if (!initTimeMs) {
      this.contactTime = new Date().getTime();
    } else {
      this.contactTime = initTimeMs;
    }
  }
  getContactTime(): number {
    return this.contactTime;
  }
  isRecentContact(withinLastMs: number): boolean {
    return (new Date().getTime() - this.contactTime) < withinLastMs;
  }
}


class LiveStreamFeedController {

  httpSubscription: Subscription | undefined;
  observers: any[] = []

  private outstandingTimer?: any;
  private lastReceivedBlockHeight = 0;

  constructor(private config: AppConfigService) {
  }

  getNextPollInternalMs() {
    return this.config.getInitialLiveStreamPollInternalMs();
  }

  clearOutstandingTimer(): void {
    if (!!this.outstandingTimer) {
      clearTimeout(this.outstandingTimer);
    }
    this.outstandingTimer = undefined;
  }

  startNextTimer(action: () => void) {
    this.outstandingTimer = setTimeout(action, this.getNextPollInternalMs());
  }

  dataReceived(payload: PurchaseListedEventPayload): void {
    // Determine max height so
    // - don't send duplicate events to our observers
    // - reduce HTTP payloads
    payload.Data.forEach(e => {
      if (e.BlockHeight > this.lastReceivedBlockHeight) {
        this.lastReceivedBlockHeight = e.BlockHeight;
      }
    });
  }

  getNextUrl(): string {
    if (this.lastReceivedBlockHeight == 0) {
      return this.config.getLastestMomentEventsStreamUrl();
    } else {
      return this.config.getMomentEventsStreamFromBlockHeightUrl(this.lastReceivedBlockHeight);
    }
  }

}


export enum ConnectionStatus {
  Unknown,
  Connected,
  Disconnected
}

@Injectable({
  providedIn: 'root'
})
export class EventBusService {

  public ConnectionStatus: BehaviorSubject<ConnectionStatus>
  public SuccessfulCollectorContact: BehaviorSubject<ContactedEvent>
  public LiveStream: Observable<PurchaseListedEvent>;
  public DetectedMoments = new EventEmitter<DetectedConsiderEvent>();

  private liveStreamController: LiveStreamFeedController;

  constructor(
    private http: HttpClient,
    private config: AppConfigService
  ) {
    this.ConnectionStatus = new BehaviorSubject<ConnectionStatus>(ConnectionStatus.Unknown);
    this.SuccessfulCollectorContact = new BehaviorSubject<ContactedEvent>(new ContactedEvent(0));
    this.liveStreamController = new LiveStreamFeedController(config);
    this.LiveStream = new Observable<PurchaseListedEvent>(this.createLiveStreamObserable());
  }

  private createLiveStreamObserable(): (this: Observable<PurchaseListedEvent>, subscriber: Subscriber<PurchaseListedEvent>) => TeardownLogic {
    const self = this;
    return (subscriber) => {
      self.liveStreamController.observers.push(subscriber);
      // When this is the first subscription, start the sequence
      if (self.liveStreamController.observers.length === 1) {
        self.startHttpSubscription();
      }
      return {
        unsubscribe() {
          // Remove from the observers array so it's no longer notified
          self.liveStreamController.observers.splice(self.liveStreamController.observers.indexOf(subscriber), 1);
          // If there's no more listeners, do cleanup
          if (self.liveStreamController.observers.length === 0) {
            self.endHttpSubscription();
          }
        }
      };
    }
  }

  private startHttpSubscription(): void {
    this.liveStreamController.httpSubscription = this.http.get(this.liveStreamController.getNextUrl()).subscribe(
      data => {
        console.log("received data");
        const payload = <PurchaseListedEventPayload> data;
        this.liveStreamController.dataReceived(payload);
        payload.Data.forEach(e => {
          this.liveStreamController.observers.forEach(s => s.next(e));
        });
      },
      error => {
        this.liveStreamController.httpSubscription = undefined;
        // TODO, have some sort of retry concept
        if (error.error instanceof ErrorEvent) {
          // client-side error
          console.log(`Client Error: ${error.error.message}`);
        } else {
          // server-side error
          console.log(`Server Error Code: ${error.status}\nMessage: ${error.message}`);
        }
        const closeNotify = this.liveStreamController.observers;
        this.liveStreamController.observers = [];
        closeNotify.forEach(s => s.error(error));
      },
      () => {
        this.liveStreamController.httpSubscription = undefined;
        this.liveStreamController.startNextTimer(() => {
          this.startHttpSubscription();
        });
      }
    );
  }

  private endHttpSubscription(): void {
    this.liveStreamController.clearOutstandingTimer();
    this.liveStreamController.httpSubscription?.unsubscribe();
    this.liveStreamController.httpSubscription = undefined;
  }

}
