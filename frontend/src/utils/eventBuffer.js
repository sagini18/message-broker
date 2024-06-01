class EventBuffer {
  constructor(dispatch, processFn, interval = 1000) {
    this.buffer = [];
    this.dispatch = dispatch;
    this.processFn = processFn;
    this.interval = interval;
    this.timer = null;
  }

  addEvent(event) {
    this.buffer.push(event);
    if (!this.timer) {
      this.startTimer();
    }
  }

  startTimer() {
    this.timer = setTimeout(() => this.flush(), this.interval);
  }

  flush() {
    if (this.buffer.length > 0) {
      const eventsToProcess = this.buffer;
      this.buffer = [];
      this.processFn(this.dispatch, eventsToProcess);
    }
    this.timer = null;
  }
}

export default EventBuffer;
