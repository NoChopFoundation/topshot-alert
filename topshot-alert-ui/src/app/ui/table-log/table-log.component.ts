import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Subscription } from 'rxjs';
import { IPlay } from 'src/app/srv/api/plays-database.service';
import { ISet } from 'src/app/srv/api/sets-database.service';
import { StaticData, StaticDataService } from 'src/app/srv/api/static-data.service';
import { Util_GetSetNameWithEditionInfo_ViaSetId } from 'src/app/srv/api/util';
import { DetectedConsiderEvent, EventBusService } from 'src/app/srv/events/event-bus.service';

interface DecoratedDetectedConsiderEvent {
  play?: IPlay;
  set?: ISet;
  ConsiderEvent: DetectedConsiderEvent;

  FullName: string;
  Type: string;
  SetName: string;
  SetDesc: string;
  SerialNumber: number;
  Price: number;
  LogLevel: string;
  DateAsLong: number;
  DateDisplay: string;
}

const ELEMENT_DATA: DecoratedDetectedConsiderEvent[] = [];

const initialSelection: DecoratedDetectedConsiderEvent[] = [];
const allowMultiSelect = true;

@Component({
  selector: 'app-table-log',
  templateUrl: './table-log.component.html',
  styleUrls: ['./table-log.component.scss'],
})

export class TableLogComponent implements OnInit, AfterViewInit, OnDestroy {

  displayedColumns: string[] = ['Type', 'FullName', 'SetName', 'SerialNumber', 'Price', 'LogLevel', 'DateAsLong', 'notes'];
  dataSource = new MatTableDataSource<DecoratedDetectedConsiderEvent>(ELEMENT_DATA);

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  selection = new SelectionModel<DecoratedDetectedConsiderEvent>(allowMultiSelect, initialSelection);

  isLoading = true;
  staticData!: StaticData;

  eventsProcessedCount = 0;

  subscriptions: Subscription[] = [];

  isAppendsFrozen = false;
  frozenDeferals: DecoratedDetectedConsiderEvent[] = [];

  constructor(
    private staticDataSrv: StaticDataService,
    private bus: EventBusService) {
  }

  ngOnInit(): void {
    this.isLoading = true;
    this.staticDataSrv.data$.subscribe(data => {
      this.staticData = data;
      this.subscriptions.push(this.bus.DetectedMoments.subscribe(evt => {
        const play = this.staticData.findPlay(evt.Event.PlayId);
        const set = this.staticData.findSet(evt.Event.SetId);
        const now = new Date();

        let SetDesc;
        if (play && set) {
          SetDesc = Util_GetSetNameWithEditionInfo_ViaSetId(this.staticData.sets, play, set.SetId);
        } else {
          SetDesc = "N/A("+evt.Event.SetId+")"
        }

        this.onEvent(<DecoratedDetectedConsiderEvent> {
          ConsiderEvent: evt,
          play: play,
          set: set,
          FullName: play ? play.FullName : "N/A("+evt.Event.PlayId+")",
          Type: evt.Event.Type,
          SetName: set ? set.Name :  "N/A("+evt.Event.SetId+")",
          SetDesc,
          SerialNumber: evt.Event.SerialNumber,
          Price: evt.Event.Price,
          LogLevel: evt.LogLevel,
          DateAsLong: now.getTime(),
          DateDisplay: now.toLocaleString()
        });
      }));
      this.subscriptions.push(this.bus.LiveStream.subscribe(evt => {
        this.eventsProcessedCount ++;
      }));
      this.isLoading = false;
    });
  }

  ngOnDestroy(): void {
    this.subscriptions.forEach(s => s.unsubscribe());
  }

  onTest() {
    console.log(this.selection.selected);
    //this.dataSource.filter = "hydrogen";
    /**
    this.dataSource.data.push({
      name: "added",
      position: 5,
      symbol: "ad",
      weight: 7
    });
    this.dataSource._updateChangeSubscription();
 */
  }

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;
  }

  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.data.length;
    return numSelected == numRows;
  }

  /** Selects all rows if they are not all selected; otherwise clear selection. */
  masterToggle() {
    this.isAllSelected() ?
        this.selection.clear() :
        this.dataSource.data.forEach(row => this.selection.select(row));
  }

  private onEvent(evt: DecoratedDetectedConsiderEvent): void {
    if (this.isAppendsFrozen) {
      this.frozenDeferals.push(evt);
    } else {
      if (this.frozenDeferals.length > 0) {
        this.frozenDeferals.forEach(d => this.dataSource.data.unshift(d));
        this.frozenDeferals = [];
      }
      this.dataSource.data.unshift(evt);
      this.dataSource._updateChangeSubscription();
    }
  }

  clearLog(): void {
    if (confirm("Delete all log entries?")) {
      this.dataSource.data = [];
      this.dataSource._updateChangeSubscription();
    }
  }

}
