import { TestBed, waitForAsync } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { AppComponent } from './app.component';
import { PlaysDatabaseTestingModule } from './srv/api/plays-database.service.spec';

describe('AppComponent', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        RouterTestingModule, PlaysDatabaseTestingModule
      ],
      declarations: [
        AppComponent
      ],
    }).compileComponents();
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });

  /**
  it(`should have as title 'topshot-alert'`, () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app.title).toEqual('topshot-alert');
  });

  it('should render title', waitForAsync (() => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;


    fixture.detectChanges();
    const compiled = fixture.nativeElement;
    expect(compiled.querySelector('.content span').textContent).toContain('topshot-alert app is running!');

    fixture.whenStable().then(() => {
      fixture.detectChanges();
      expect(fixture.nativeElement.querySelector('.content span').textContent).toContain('change it up app is running2!');
    });
    app.ngOnInit();
    fixture.autoDetectChanges();
    app.title = "change it up";

  }));
   */
});
