class AddStatsToPlayers < ActiveRecord::Migration[5.2]
  def change
    # Stolen from https://en.wikipedia.org/wiki/Attribute_(role-playing_games)
    #
    # These were chosen over stats like "Defense" because their purpose is
    # to steer character choices in the story, and not to judge determine the
    # outcomes in a fight.

    add_column :players, :strength,  :integer, null: false, default: 0
    add_column :players, :charisma,  :integer, null: false, default: 0
    add_column :players, :intellect, :integer, null: false, default: 0
    add_column :players, :agility,   :integer, null: false, default: 0
    add_column :players, :luck,      :integer, null: false, default: 0

    # Level
    add_column :players, :level, :integer, null: false, default: 1
  end
end
