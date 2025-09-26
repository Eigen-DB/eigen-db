import { Heading, Highlight } from "@chakra-ui/react";

export default function Header() {
  return (
    <Heading size="3xl" style={{ textAlign: 'center'}}>
      <Highlight
        query="meaning"
        styles={{ px: "0.5", bg: "purple.subtle" }}
      >
        Search by meaning, not just words.
      </Highlight>
    </Heading>
  );
}