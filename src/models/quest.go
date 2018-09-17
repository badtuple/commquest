package models

// A Quest's narrative consists of 4 sections:
//	Intro -- a Call to action for the group.
//		Examples:
//			a villager needs an item
//			a town is being terrorized
//			a person needs rescued
//
//	Journey -- an array of obstacles that need to be overcome.
//		Each array element takes 1 turn to complete.
//		Examples:
//			a dungeon must be explored
//			a person must be defeated
//			a person must be found
//
//	Return -- Completely optional step where an additional
//		Journey element must be passed through on the
//		return to the origin of the quest.
//
//	End -- The original Call to action is confronted and
//		the characters are rewarded with XP, items, etc.
//		Examples:
//			villager accepts needed item
//			town is grateful that you've saved them
//			person is safely returned home
type Quest struct {
	// Narrative structure
	Intro   intro
	Journey []obstacle
	Return  *obstacle
	End     end
}

type intro struct {
	// Possible intro types:
	//	obtain    -- Obtain an item or animal
	//	destroy   -- Destroy an item
	//	transport -- Transport an item
	//	save_town -- Save a town that's being terroized
	//	rescue    -- Rescue a person or animal
	//	escort    -- Escort a person/animal between towns
	Type string

	Tmpl string // Intro text template

	// Item to be obtained, destroyed, or transported
	Item *string

	// Person to be rescued or escorted
	Person *string

	// Town to be saved. Alternatively the town where
	// the quest originated.
	Town *string

	// Certain intro types require a specific ending
	EndType []string
}

type obstacle struct{}

type end struct {
	// Possible end types:
	//	none   -- Quest just ends.
	// 	thanks -- The person who sent you on the quest says thank you.
	//	party  -- The town celebrates their safety.
	Type string
	Tmpl string // Ending text template
}

// TODOs:
//	[ ] obtain -- Obtain a particular item or animal
//	[ ] destroy -- Destroy a particular item
//	[ ] transport -- Transport a particular item
//	[ ] save_town -- Save a town that's being terroized
//	[ ] rescue -- Rescue a particular person or animal
//	[ ] escort -- Transport a person or animal between towns
