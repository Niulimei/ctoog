/* eslint no-useless-escape:0 import/prefer-default-export:0 */
const reg = /(((^https?:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+(?::\d+)?|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)$/;

export const isUrl = (path: string): boolean => reg.test(path);

export function guid() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    // eslint-disable-next-line
    let r = (Math.random() * 16) | 0,
      // eslint-disable-next-line
      v = c == 'x' ? r : (r & 0x3) | 0x8;

    return v.toString(16);
  });
}

const parseDuration = (duration: number) => {
  const hours: number = Math.floor(duration / 3600);
  const minutes: number = Math.floor((duration - hours * 3600) / 60);
  const seconds: number = Math.floor(duration - hours * 3600 - minutes * 60);

  return {
    hours,
    minutes,
    seconds,
  };
};

export const formatDuration = (duration: number) => {
  const fillZero = (num: number) => num.toString().padStart(2, '0');
  const { hours, minutes, seconds } = parseDuration(duration);

  return `${fillZero(hours)}:${fillZero(minutes)}:${fillZero(seconds)}`;
};

export function humanizeDuration(duration: number) {
  const { hours, minutes, seconds } = parseDuration(duration);

  // eslint-disable-next-line no-nested-ternary
  return hours ? `${hours} h` : minutes ? `${minutes} min` : `${seconds} sec`;
}
