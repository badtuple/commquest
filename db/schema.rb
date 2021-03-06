# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 2018_09_17_123526) do

  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "items", force: :cascade do |t|
    t.string "name", null: false
    t.string "article", null: false
    t.integer "xp_incr", default: 0, null: false
    t.integer "level_incr", default: 0, null: false
    t.integer "strength_incr", default: 0, null: false
    t.integer "charisma_incr", default: 0, null: false
    t.integer "intellect_incr", default: 0, null: false
    t.integer "agility_incr", default: 0, null: false
    t.integer "luck_incr", default: 0, null: false
  end

  create_table "players", force: :cascade do |t|
    t.string "handle", null: false
    t.string "name", null: false
    t.string "class", null: false
    t.integer "xp", default: 0, null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.integer "strength", default: 0, null: false
    t.integer "charisma", default: 0, null: false
    t.integer "intellect", default: 0, null: false
    t.integer "agility", default: 0, null: false
    t.integer "luck", default: 0, null: false
    t.integer "level", default: 1, null: false
  end

  create_table "quests", force: :cascade do |t|
    t.string "state", null: false
    t.boolean "inprogress", default: false, null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

end
