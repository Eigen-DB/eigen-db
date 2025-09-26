"use client";

import { Card, Button, Badge } from "@chakra-ui/react";
import { HiAtSymbol } from "react-icons/hi";
import { FaSearch } from "react-icons/fa"

export default function Tweet({ user, body, getSimilarHandler }: { user: string; body: string; getSimilarHandler: (tweet: string) => void }) {
  return (
    <Card.Root width="70%" minWidth={"500px"} maxWidth={"600px"}>
      <Card.Body gap="2">
        <Card.Title mt="2">
          <Badge colorPalette="teal" size="md" variant="subtle" style={{borderRadius: '25px'}}>
            <HiAtSymbol/>{user}
          </Badge>
        </Card.Title>
        <Card.Description>
          {body}
        </Card.Description>
      </Card.Body>
      <Card.Footer justifyContent="flex-end">
        <Button variant={"surface"} colorPalette={"purple"} onClick={() => getSimilarHandler(body)}><FaSearch/> Explore similar</Button>
      </Card.Footer>
    </Card.Root>
  );
}
