import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QuickGraphComponent } from './quick-graph.component';

describe('QuickGraphComponent', () => {
  let component: QuickGraphComponent;
  let fixture: ComponentFixture<QuickGraphComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ QuickGraphComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(QuickGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
