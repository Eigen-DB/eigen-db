import Header from "./components/header";
import SearchParent from "./components/searchParent";

export default function Home() {
  let tweets = [
    {
      id: "1",
      user: "Young_J",
      content: "I'm off too bed. I gotta wake up hella early tomorrow morning."
    },
    {
      id: "2",
      user: "szrhnds602",
      content: "Borders closed at 10"
    },
    {
      id: "3",
      user: "MichaelPe",
      content: "@FollowSavvy I never found her. everytime I click on her twitter thing through your myspace..... it goes to some dude's page"
    },
    {
      id: "4",
      user: "kissbliss1",
      content: "18 weeks till sisters home.. i missed her call, again! its the worst feeling in the world."
    },
    {
      id: "5",
      user: "FuhQ",
      content: "Car show season has started without me"
    }
  ]

  return (
    <div style={{ margin: "40px auto" }}>
      <Header />
      <SearchParent tweets={tweets} />
      <footer style={{bottom: 0, left: 0, width: "100%", textAlign: "center", padding: "25px", color: 'grey', fontSize: '14px'}}>
        Built using 1000 rows from <a href="https://www.kaggle.com/datasets/kazanova/sentiment140/data" target="_blank" style={{fontWeight: 'bold'}}><u>this</u></a> dataset :)
      </footer>
    </div>
  );
}

