import { colors } from "../../settings/theme";
import classes from "../../utils/classes";
import { StepperProps } from "./types";
import style from "./style/Stepper.module.scss";
import { PiCheckBold } from "react-icons/pi";

export default function Stepper({ steps, current, clickable, setCurrent }: StepperProps) {
  return (
    <div className={style.stepper}>
      {steps.map((step, i) => {
        const done = i < current;

        const Element = clickable ? "button" : "div";

        const buttonProps: any = clickable
          ? {
              onClick: () => setCurrent(i),
              type: "button",
            }
          : {};

        return (
          <Element
            key={i}
            className={classes([style.step, i === current ? style.active : done && style.done, steps.length < 4 && style.stepWide])}
            {...buttonProps}>
            <div className={style.badge}>{done ? <PiCheckBold color={colors.primary} /> : i + 1}</div>
            <div className={style.label}>{step.label}</div>
          </Element>
        );
      })}
    </div>
  );
}
