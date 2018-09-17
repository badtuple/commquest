class CreateQuest < ActiveRecord::Migration[5.2]
  def change
    create_table :quests do |t|
      # The current serialized state of the Quest
      t.string  :state,      null: false
      t.boolean :inprogress, null: false, default: false
      t.timestamps
    end
  end
end
