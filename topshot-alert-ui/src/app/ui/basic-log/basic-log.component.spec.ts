import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BasicLogComponent } from './basic-log.component';

describe('BasicLogComponent', () => {
  let component: BasicLogComponent;
  let fixture: ComponentFixture<BasicLogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ BasicLogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(BasicLogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
