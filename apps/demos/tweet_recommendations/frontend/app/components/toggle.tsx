"use client";

import { Switch } from "@chakra-ui/react";

export default function Toggle({ onToggle }: { onToggle: (newState: boolean) => void }) {
  return (
    <Switch.Root colorPalette={"purple"} defaultChecked onCheckedChange={({checked}) => onToggle(checked)} >
      <Switch.HiddenInput />
      <Switch.Control />
      <Switch.Label>Use vector similarity</Switch.Label>
    </Switch.Root>
  )
}
