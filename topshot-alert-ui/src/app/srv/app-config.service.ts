import { Injectable } from '@angular/core';
import { environment } from '../.././environments/environment';
import { IPlay } from './api/plays-database.service';
import { ISet } from './api/sets-database.service';
import { Util_GetSetNameWithEditionInfo_ViaSetId } from './api/util';

declare function require(name:string): any;
const jsStringEscape = require('js-string-escape');

export interface IAppConfig {
  production: boolean
}

@Injectable({
  providedIn: 'root'
})
export class AppConfigService {

  private config: IAppConfig

  constructor() {
    this.config = environment
  }

  isProduction(): boolean {
    return this.config.production
  }

  getPlaysDataUrl(): string {
    return this.getPlayDataEndpoint() + "/plays";
  }

  getSetsDataUrl(): string {
    return this.getPlayDataEndpoint() + "/sets";
  }

  getRecentCollectorsUrl(): string {
    return this.getPlayDataEndpoint() + "/status/collectors/recent";
  }

  getLastestMomentEventsStreamUrl(): string {
    return this.getPlayDataEndpoint() + "/stream/momentEvents/from/HEAD";
  }

  getMomentEventsStreamFromBlockHeightUrl(lastBlockHeight: number): string {
    return this.getPlayDataEndpoint() + `/stream/momentEvents/from/HEAD/to/blockHeight/${lastBlockHeight}`;
  }

  getRawLinkPlaySetUrl(playId: number, setId: number): string {
    return this.getPlayDataEndpoint() + `/momentEvents/play/${playId}/set/${setId}/from/HEAD`;
  }

  getCollectorPingRetryMs(): number {
    return 60*1000;
  }

  getInitialLiveStreamPollInternalMs() {
    return 15*1000;
  }

  getInitialConsiderScript(plays: IPlay[], sets: ISet[], playId: number, setId: number) {
    const play = this.findPlay(plays, playId);
    const set = this.findSet(sets, setId);
    const momentDescription = jsStringEscape("Note for " + play.FullName + " " + Util_GetSetNameWithEditionInfo_ViaSetId(sets, play, setId));
    return "function consider(API, evt) {\n" +
      `  // ${play.FullName}(${playId}) - ${play.PlayType} - ${set.Name}(${setId}) \n` +
      `  if ((evt.isListing() || evt.isPurchase()) && evt.PlayId === ${playId} && evt.SetId === ${setId}) {\n` +
      "    if (evt.Price < 10 && evt.SerialNumber < 5000) {\n" +
      `      API.log(API.LogLevel.ALERT, evt, "${momentDescription}");\n` +
      "    } else {\n" +
      `      API.log(API.LogLevel.CONSOLE, evt, "${momentDescription}");\n` +
      "    }\n" +
      "  } else {\n" +
      "    // API.log(API.LogLevel.DEBUG, evt);\n" +
      "  }\n" +
      "}\n";
  }

  private getPlayDataEndpoint(): string {
    TODO POINT TO YOUR BACKEND
    return "https://************.cloudfunctions.net/GetPlayDataHTTP"
  }

  private findPlay(plays: IPlay[], playId: number): IPlay {
    return plays.find(e => e.PlayId === playId)!;
  }

  private findSet(sets: ISet[], setId: number): ISet {
    return sets.find(e => e.SetId === setId)!;
  }

}
