import { HttpClient } from '@angular/common/http';
import { AfterViewInit, Component, Input, OnChanges, OnInit, SimpleChanges, ViewChild } from '@angular/core';
import { ChartType, Column, GoogleChartComponent } from 'angular-google-charts';
import { AppConfigService } from 'src/app/srv/app-config.service';
import { PurchaseListedEvent, PurchaseListedEventPayload } from 'src/app/srv/events/event-bus.service';

const TOOLTIP_COLUMN = {type: "string", role: "tooltip"};

@Component({
  selector: 'app-quick-graph',
  templateUrl: './quick-graph.component.html',
  styleUrls: ['./quick-graph.component.scss']
})
export class QuickGraphComponent implements OnInit, AfterViewInit, OnChanges {

  isLoading = true;
  isGraphable = false;
  serialNumberGuid = "SerialNumber Guide - diamond(2-digit) square(3-digit) triangle(< 2K) circle(> 2K)";

  private queryResult: PurchaseListedEvent[] = [];

  @Input() playId?: number;
  @Input() setId?: number;

  @ViewChild('chart', {static: false}) scatterChart!: GoogleChartComponent;

  data:any = [];
  type: ChartType = ChartType.ScatterChart;
  columns: Column[] = [];
  options: any = {};

  width = 600;
  height = 400;

  constructor(
    private config: AppConfigService,
    private http: HttpClient
  ) { }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.playId || changes.setId) {
      if (this.playId && this.setId) {
        this.updateChart();
      }
    }
  }

  ngAfterViewInit(): void {
  }

  ngOnInit(): void {
  }

  private updateChart(): void {
    this.isLoading = true;
    this.http.get<PurchaseListedEventPayload>(this.config.getRawLinkPlaySetUrl(this.playId!, this.setId!)).subscribe(payload => {
      this.queryResult = payload.Data;
      const filteredEvts = this.queryResult.filter(evt => evt.Type === "P");
      if (filteredEvts.length === 0) {
        this.isGraphable = false;
      } else {
        this.isGraphable = true;
        this.setupDefaultChart(filteredEvts);
      }
      this.isLoading = false;
    });
  }

  private setupDefaultChart(filteredEvts: PurchaseListedEvent[]): void {
    let title;
    let nData: any[] = [];
    let maxValueX = 0;
    let minValueX = 0;
    let xLabel;

    title = "Purchases";
    const uColumns: Column[] = [ "Time" ];

    const stratified: any[] = [];
    let seriesIndicators: { [key: number]: any} = {};
    let seriesIndicatorIndex = 0;

    const gLow = filteredEvts.filter(evt => evt.SerialNumber < 100);
    if (gLow.length > 0) {
      uColumns.push("2-digit");
      uColumns.push(TOOLTIP_COLUMN);
      stratified.push(gLow);
      seriesIndicators[seriesIndicatorIndex++] = { pointShape: 'diamond' };
    }

    const gMid = filteredEvts.filter(evt => evt.SerialNumber >= 100 && evt.SerialNumber < 1000);
    if (gMid.length > 0) {
      uColumns.push("3-digit");
      uColumns.push(TOOLTIP_COLUMN);
      stratified.push(gMid);
      seriesIndicators[seriesIndicatorIndex++] = { pointShape: 'square' };
    }

    const gLarge = filteredEvts.filter(evt => evt.SerialNumber >= 1000 && evt.SerialNumber < 2000);
    if (gLarge.length > 0) {
      uColumns.push("Less 2K");
      uColumns.push(TOOLTIP_COLUMN);
      stratified.push(gLarge);
      seriesIndicators[seriesIndicatorIndex++] = { pointShape: 'triangle' };
    }

    const gRemain = filteredEvts.filter(evt => evt.SerialNumber >= 2000);
    if (gRemain.length > 0) {
      uColumns.push("High");
      uColumns.push(TOOLTIP_COLUMN);
      stratified.push(gRemain);
      seriesIndicators[seriesIndicatorIndex++] = { pointShape: 'circle' };
    }

    // Ordered by Created_At desc
    const lastIdx = filteredEvts.length - 1;
    const maxMs = new Date().getTime();
    minValueX = -1 * this.convertRelativeHours(filteredEvts[lastIdx], maxMs);

    for (var stratLevel = 0; stratLevel < stratified.length; stratLevel++) {
      stratified[stratLevel].forEach((evt: PurchaseListedEvent) => {
        const row: any[] = [];
        row.push(-1 * this.convertRelativeHours(evt, maxMs)); // Time column

        // Add nulls before
        for (var i = 0; i < stratLevel; i++) {
          row.push(null);
          row.push(null);
        }

        row.push(evt.Price);
        row.push(`SN(${evt.SerialNumber}) / $${evt.Price}`);

        // Add nulls after
        for (var i = stratLevel+1; i < stratified.length; i++) {
          row.push(null);
          row.push(null);
        }
        nData.push(row);
      });
    }

    xLabel = "Hours " + filteredEvts[lastIdx].Created_At + " to " + filteredEvts[0].Created_At;

    this.columns = uColumns;
    this.data = nData;
    this.options = {
      title,
      hAxis: {title: xLabel, minValue: minValueX, maxValue: maxValueX},
      vAxis: {title: 'Price'},
      legend: 'none',
      series: seriesIndicators
    };
  }

  private convertRelativeHours(pEvt: PurchaseListedEvent, maxMs: number): number {
    return (maxMs - Date.parse(pEvt.Created_At)) / (1000*60*60);
  }

}
