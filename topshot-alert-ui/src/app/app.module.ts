import { CUSTOM_ELEMENTS_SCHEMA, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { MatTableModule } from '@angular/material/table';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { CollectorStatusDetectorComponent } from './api/collector-status-detector/collector-status-detector.component';
import { PlayerLookupComponent } from './ui/player-lookup/player-lookup.component';
import { PlayerLookupResultListComponent } from './ui/player-lookup-result-list/player-lookup-result-list.component';
import { PlayerMonitorComponent } from './ui/player-monitor/player-monitor.component';
import { PlayerMonitorItemComponent } from './ui/player-monitor-item/player-monitor-item.component';
import { PlayerMonitorItemEditorComponent } from './ui/player-monitor-item-editor/player-monitor-item-editor.component';

import { FormsModule } from '@angular/forms';
import { CodemirrorModule } from '@ctrl/ngx-codemirror';
import { BasicLogComponent } from './ui/basic-log/basic-log.component';
import { TableLogComponent } from './ui/table-log/table-log.component';
import { MatSortModule } from '@angular/material/sort';
import { MatPaginatorModule} from '@angular/material/paginator';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { QuickGraphComponent } from './ui/quick-graph/quick-graph.component';
import { GoogleChartsModule } from 'angular-google-charts';

@NgModule({
  declarations: [
    AppComponent,
    CollectorStatusDetectorComponent,
    PlayerLookupComponent,
    PlayerLookupResultListComponent,
    PlayerMonitorComponent,
    PlayerMonitorItemComponent,
    PlayerMonitorItemEditorComponent,
    BasicLogComponent,
    TableLogComponent,
    QuickGraphComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    FormsModule,
    CodemirrorModule,
    MatTableModule,
    MatPaginatorModule,
    MatSortModule,
    MatCheckboxModule,
    BrowserAnimationsModule,
//    GoogleChartsModule,
    GoogleChartsModule.forRoot()
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }

