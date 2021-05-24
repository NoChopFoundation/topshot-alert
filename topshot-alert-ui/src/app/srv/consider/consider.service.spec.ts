import { TestBed } from '@angular/core/testing';
import { PurchaseListedEvent } from '../events/event-bus.service';

import { ConsiderSandbox, ConsiderService } from './consider.service';

const SCRIPT1 = "function consider(API, evt) {\n" +
"  // James Harden(338) - Rim - Platinum Ice(3) \n" +
"  if (evt.isListing() && evt.PlayId === 338 && evt.SetId === 3) {\n" +
"    if (evt.Price < 10 && evt.SerialNumber < 5000) {\n" +
"      API.log(API.LogLevel.ALERT, evt, \"Cheap Harden\");\n" +
"    } else {\n" +
"      API.log(API.LogLevel.CONSOLE, evt);\n" +
"    }\n" +
"  }\n" +
"}\n";

const SCRIPT2 = `(

  function sandbox(api, evt) {
    function consider(API, evt) {
      // James Harden(157) - 3 Pointer - Base Set(2)
      if (evt.isListing() && evt.PlayId === 157 && evt.SetId === 2) {
        if (evt.Price < 10 && evt.SerialNumber < 5000) {
          API.log(API.LogLevel.ALERT, evt, "Cheap Harden");
        } else {
          API.log(API.LogLevel.CONSOLE, evt);
        }
        }
    }
    consider(api, evt);
  }

  )`;

describe('ConsiderService', () => {
  let service: ConsiderService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ConsiderService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('basic sandbox', () => {
    debugger;
    const sb = new ConsiderSandbox();
    const evt = <PurchaseListedEvent> {
      BlockHeight: 999,
      MomentId: 898,
      PlayId: 338,
      Price: 5,
      SetId: 3,
      SellerAddr: "addrA",
      SerialNumber: 44,
      Type: "P"
    };

    const entryFunc = eval(sb.compileUserCode(SCRIPT2));
    debugger;
    entryFunc(sb.createApi(), sb.createConsiderEvent(evt));
    debugger;

  });
});

