"use client";

import Tweet from "./tweet";
import Toggle from "./toggle";
import SearchBar from "./searchBar";
import { useState } from "react";
import { Spinner } from "@chakra-ui/react";

export default function SearchParent({ tweets }: { tweets: { id: string, user: string; content: string }[] }) {
  const [useSimilarity, setUseSimilarity] = useState(true);
  const [tweetsState, setTweetsState] = useState(tweets);
  const [searching, setSearching] = useState(false);

  const getSimilarHandler = async(tweet: string) => {
    setSearching(true);
    setTweetsState([]);
    const response = await fetch(
      'http://localhost:5000/recommend', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({tweet: tweet, use_eigendb: useSimilarity, desired_results: 5})
      }
    );

    if (!response.ok) {
      alert('Failed to fetch recommendations');
      return;
    }

    const data = await response.json();
    const tweets: { id: string, user: string; content: string }[] = [];
    for (const rec of data.recommendations) {
      tweets.push(rec.metadata);
    }
    setSearching(false);
    setTweetsState(tweets);
  }

  return (
    <>
      <div style={{ margin: "24px 0" }}>
        <Toggle
          onToggle={setUseSimilarity}
        />
      </div>

      <SearchBar
        onSearch={getSimilarHandler}
      />

      {searching ? <Spinner size="md" style={{marginTop: '20px'}} /> : null}
      
      <div style={{ display: "flex", flexDirection: "column", gap: "20px", marginTop: '20px', alignItems: "center" }}>
        {tweetsState.length > 0 ? tweetsState.map(({id, user, content}) => (
          <Tweet 
            key={id} 
            user={user} 
            body={content}
            getSimilarHandler={getSimilarHandler}
          />
        )) : <p>{searching ? 'Searching...' : 'No tweets found.'}</p>}
      </div>
    </>
  );
} 