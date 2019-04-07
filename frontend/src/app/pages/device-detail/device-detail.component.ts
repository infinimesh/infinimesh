import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {DeviceService} from "../../core/data/device.service";

@Component({
  selector: 'device-detail',
  templateUrl: './device-detail.component.html',
  styleUrls: ['./device-detail.component.scss']
})
export class DeviceDetailComponent implements OnInit, OnDestroy {

  device$;
  state;
  stateSubscription;
  desiredStateUpdated = false;
  reportedStateUpdated = false;
  JSON = JSON;

  constructor(private route: ActivatedRoute,
              private deviceService: DeviceService) { }

  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      const deviceId = params.get("deviceId")
      this.device$ = this.deviceService.getOne(deviceId);
      this.deviceService.getState(deviceId).subscribe(state => {
        if(state) {
          this.state = state;
        } else {
          this.state = {
            reported: {},
            desired: {}
          };
        }
        this.stateSubscription = this.deviceService.streamState(deviceId).subscribe((data) => {
          this.handleStateUpdate(data);
        });
      })
    })
  }

  private handleStateUpdate(data) {
    if (data.reportedState && data.reportedState !== null) {
      this.state.reported = data.reportedState;
      this.reportedStateUpdated = true
      setTimeout(() => {
        this.reportedStateUpdated = false
      }, 500);
    }
    if (data.desiredState && data.desiredState !== null) {
      this.state.desired = data.desiredState;
      this.desiredStateUpdated = true
      setTimeout(() => {
        this.desiredStateUpdated = false
      }, 500);
    }
  }

  ngOnDestroy(): void {
    this.stateSubscription.unsubscribe();
  }

}
