"use client";

import { Button, Input, Spinner } from "@chakra-ui/react"
import { FaSearch } from "react-icons/fa"
import { useState } from "react";
import { Badge, Stack } from "@chakra-ui/react"

export default function SearchBar({ onSearch }: { onSearch: (tweet: string) => void }) {
  const [inputValue, setInputValue] = useState("");
  return (
    <>  
        <div>
          <Input 
            placeholder="Search tweets by meaning..." 
            value={inputValue || ''}
            size="md" 
            width="400px"
            onChange={(e) => setInputValue(e.target.value)}
            style={{ borderTopRightRadius: 0, borderBottomRightRadius: 0 }}
          />
          <Button 
            variant={"subtle"}
            style={{ borderTopLeftRadius: 0, borderBottomLeftRadius: 0 }}
            onClick={() => onSearch(inputValue)}
            disabled={inputValue.length === 0}
          >
            <FaSearch/>Find tweets
          </Button>
        </div>

        <Stack direction="row" mt="10px" style={{ justifyContent: "center" }}>
          <Badge colorPalette="purple" size="sm" variant="surface" cursor="pointer" onClick={() => { setInputValue("travel"); onSearch("travel") }}><FaSearch/> travel</Badge>
          <Badge colorPalette="purple" size="sm" variant="surface" cursor="pointer" onClick={() => { setInputValue("concert"); onSearch("concert") }}><FaSearch/> concert</Badge>
          <Badge colorPalette="purple" size="sm" variant="surface" cursor="pointer" onClick={() => { setInputValue("studying"); onSearch("studying") }}><FaSearch/> studying</Badge>
          <Badge colorPalette="purple" size="sm" variant="surface" cursor="pointer" onClick={() => { setInputValue("music"); onSearch("music") }}><FaSearch/> music</Badge>
        </Stack>
    </>
  )
}