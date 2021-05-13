import React from 'react';
import { render } from 'react-dom';

interface Line {
  number: number;
  text: string;
}

const LineComponent: React.FC<{ text: string }> = ({ text }) => {
  return (
    <div className="line">
      <a></a>
      <span>{text}</span>
    </div>
  );
};

export default class Log {
  constructor(container: HTMLElement) {
    this.container = container;
  }

  container: HTMLElement;
  line = {} as Line;
  lines = [] as Line[];

  renderLines() {
    render(
      // eslint-disable-next-line react/no-array-index-key
      this.lines.map((line, index) => <LineComponent key={index} text={line.text} />),
      this.container,
    );
  }

  writeNewLine(text: string) {
    this.line = { number: this.line.number + 1, text };
    this.lines = this.lines.concat(this.line);
  }

  clearOutput() {
    this.lines.length = 0;
    this.renderLines();
  }

  write(values: string) {
    const regExp = /(.+?)(\r?\n|$)/g;
    let res;
    // eslint-disable-next-line no-cond-assign
    while ((res = regExp.exec(values)) != null) {
      this.writeNewLine(res[1]);
    }
    this.renderLines();
  }
}
